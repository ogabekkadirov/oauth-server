package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/config"
)

func NewPostgresPool() *pgxpool.Pool {
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

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=%d&sslmode=disable",
	config.PostgresUser,
	config.PostgresPassword,
	config.PostgresHost,
	config.PostgresPort,
	config.PostgresDatabase,
	60,)
	configdb, err := pgxpool.ParseConfig(dbURL)
    if err != nil {
        log.Fatalf("Failed to parse DSN: %v", err)
    }
	configdb.MaxConns = 10
    configdb.MinConns = 2
    configdb.MaxConnLifetime = time.Minute * 5

    pool, err := pgxpool.NewWithConfig(context.Background(), configdb)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }

    // Optional: test connection
    if err := pool.Ping(context.Background()); err != nil {
        log.Fatalf("DB ping failed: %v", err)
    }

    log.Println("✅ Connected to Postgres via pgx")
    return pool
}
