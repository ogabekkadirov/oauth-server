package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/middlewares"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/response"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/utils/cerror"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/services"
)


type authController struct {
	authSvc services.AuthService
}
func (r *router) initAuthController() {
	h := &authController{
		authSvc: r.authSvc,
	}

	r.routes.v1.POST("/oauth/token", h.handleToken)
	r.routes.v1.GET("/me",middlewares.AuthenticateMiddleware(),h.handleMe)
}
func (h *authController) handleToken(ctx *gin.Context) {
	var req models.TokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cErr := cerror.NewError(http.StatusBadRequest, err) 
		response.ErrorResult(ctx,cErr)
		return
	}

	var token *models.Token
	var err error
	switch req.GrantType {
	case "client_credentials":
		token, err = h.authSvc.HandleClientCredentials(req.ClientID, req.ClientSecret)

	case "password":
		token, err = h.authSvc.HandlePasswordGrant(req.Username, req.Password, req.ClientID,req.ClientSecret)

	case "authorization_code":
		token, err = h.authSvc.HandleAuthorizationCodeGrant(req.Code, req.ClientID, req.RedirectURL,req.ClientSecret)

	case "refresh_token":
		token, err = h.authSvc.HandleRefreshToken(req.RefreshToken, req.ClientID,req.ClientSecret)

	default:
		err = errors.New("unsupported grant_type")
		cErr := cerror.NewError(http.StatusBadRequest, err) 
		response.ErrorResult(ctx,cErr)
		return
	}

	if err != nil {
		cErr := cerror.NewError(http.StatusBadRequest, err) 
		response.ErrorResult(ctx,cErr)
		return
	}

	response.SuccessResult(ctx, http.StatusOK, token)
}

func (h *authController) handleMe(ctx *gin.Context) {
	userID, isExist := ctx.Get("AuthUserId")
	if !isExist {
		cErr := cerror.NewError(http.StatusUnauthorized, errors.New("Unauthorized")) 
		response.ErrorResult(ctx,cErr)
		return
	}

	user, err := h.authSvc.GetUserByID(userID.(string))
	if err != nil {
		cErr := cerror.NewError(http.StatusNotFound, err) 
		response.ErrorResult(ctx,cErr)
		return
	}
	response.SuccessResult(ctx, http.StatusOK, user)
}