package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/config"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/jwt"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/response"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/utils/cerror"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"
)

func AuthenticateMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var cError models.Error
		authorization := string(ctx.Request.Header.Get("Authorization"))
		if authorization != "" {
			auth := strings.SplitN(authorization, " ", 2)
			if len(auth) != 2 || auth[0] != "Bearer" {
				cError = cerror.NewError(http.StatusUnauthorized, errors.New("Unauthorized"))
				response.ErrorResult(ctx, cError)
			} else {
				config,err := config.Load()
				if err != nil{
					panic(err)
				}

				jwt, err := jwt.NewJwtService(&config)
				if err != nil{
					panic(err)
				}
				claims, err := jwt.VerifyToken(ctx,auth[1])
				if err != nil {
					cError = cerror.NewError(http.StatusUnauthorized, err)
					response.ErrorResult(ctx, cError)
				}
				ctx.Set("AuthUserId", claims.Sub)
			}

		} else {
			cError = cerror.NewError(http.StatusUnauthorized, errors.New("unauthorized!  bearer token was not found"))
			response.ErrorResult(ctx, cError)
		}
		ctx.Next()
	}
}