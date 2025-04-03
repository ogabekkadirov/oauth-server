package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/response"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/services"
)

type routes struct {
	root *gin.Engine
	v1   *gin.RouterGroup
}

type router struct {
	routes      *routes
	authSvc  services.AuthService
}

func Init(
	root *gin.Engine,
	authSvc  services.AuthService,
) {
	// ping pong
	root.GET("/ping", func(ctx *gin.Context) {
		response.SuccessResult(ctx, http.StatusOK, "pong")
	})
	v1  := root.Group("/api/v1")
	routes := routes{
		root: root,
		v1:   v1 ,
	}

	router := router{
		routes:      &routes,
		authSvc:    authSvc,
	}

	router.initAuthController()
}