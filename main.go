package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/config"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/db"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/jwt"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/redis"
	authcoderepo "github.com/ogabekkadirov/oauth-server/src/Infrastructure/repositories/auth_code"
	clientrepo "github.com/ogabekkadirov/oauth-server/src/Infrastructure/repositories/client"
	tokenrepo "github.com/ogabekkadirov/oauth-server/src/Infrastructure/repositories/token"
	userrepo "github.com/ogabekkadirov/oauth-server/src/Infrastructure/repositories/user"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/controllers"
	authsvc "github.com/ogabekkadirov/oauth-server/src/domain/auth/services"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)


func main() {
	ctx, cancel := context.WithCancel(context.Background())
	config,err := config.Load()
	if err != nil{
		panic(err)
	}
	logger, err := config.NewLogger()
	if err != nil{	
		panic(err)
	}
	defer logger.Sync()
	// db connection
	
	pgDB := db.NewPostgresPool()
	jwtGen,err := jwt.NewJwtService(&config)
	if err != nil{	
		panic(err)
	}

	// if len(os.Args) > 1 && os.Args[1] == "seed" {
	// 	seeder.RunSeeder()
	// 	return
	// }
	
	// start redis
	rdb := redis.NewRedisClient(config.RedisAddr)

	// Redisga test "ping"
	if err := rdb.RdbClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Redisga ulanishda xato: %v", err)
	}
	// end redis 

	// start repositories
	clientRepo := clientrepo.NewClientRepository(pgDB)
	userRepo := userrepo.NewUserRepository(pgDB)
	tokenStore := tokenrepo.NewTokenRepository(rdb.RdbClient)
	authCodeRepo := authcoderepo.NewAuthRepository(rdb.RdbClient)
	// end repositories

	// start services
	authService := authsvc.NewAuthService(userRepo,clientRepo, tokenStore, jwtGen,authCodeRepo)
	// end services
	
	//  start Controllers
	root := gin.Default()

	root.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowCredentials: true,
	}))

	controllers.Init(root,authService)
	// cancel context on shutdown
	g,ctx := errgroup.WithContext(ctx)

	osSignals := make(chan os.Signal,1)

	signal.Notify(osSignals,os.Interrupt,syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(osSignals)
	// stargt http server

	var httpServer *http.Server

	g.Go(func() error{
		httpServer = &http.Server{
			Addr:  ":" + config.HttpPort,
			Handler: root,
		}

		logger.Debug("main: started http server", zap.String("port",config.HttpPort))
		if err := httpServer.ListenAndServe();err != http.ErrServerClosed{
			return err
		}
		return nil
	})


	select {
	case <-osSignals:
		logger.Info("main: received os signal, shutting down")
		break
	case <-ctx.Done():
		break
	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(),5*time.Second)
	defer shutdownCancel()

	if httpServer != nil {
		httpServer.Shutdown(shutdownCtx)
	}

	if err := g.Wait(); err != nil {
		logger.Error("main: server returned an error",zap.Error(err))
	}
}
