package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/middlewares"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/response"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/utils/cerror"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/models"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/services"
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
	r.routes.v1.GET("/authorize", h.Authorize)
	r.routes.v1.POST("/login", h.Login)
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

func (h *authController) Authorize(ctx *gin.Context) {
	clientID := ctx.Query("client_id")
    redirectURI := ctx.Query("redirect_uri")
    scope := ctx.Query("scope")
    state := ctx.Query("state")
	userID, isExist := ctx.Get("AuthUserId")
    if !isExist {
        loginURL := fmt.Sprintf("/login?redirect_uri=%s&client_id=%s&scope=%s&state=%s",
            url.QueryEscape(redirectURI),
            url.QueryEscape(clientID),
            url.QueryEscape(scope),
            url.QueryEscape(state),
        )
        ctx.Redirect(http.StatusFound, loginURL)
        return
    }
	h.CompleteAuthorization(ctx, clientID, userID.(string), redirectURI, scope, state)
}

func (h *authController) Login(ctx *gin.Context) {
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
    redirectURI := ctx.PostForm("redirect_uri")
    clientID := ctx.PostForm("client_id")
    scope := ctx.PostForm("scope")
    state := ctx.PostForm("state")
    user, err := h.authSvc.ValidateUser(username, password)
    if err != nil {
		loginURL := fmt.Sprintf("/login?redirect_uri=%s&client_id=%s&scope=%s&state=%s&error=%s",
            url.QueryEscape(redirectURI),
            url.QueryEscape(clientID),
            url.QueryEscape(scope),
            url.QueryEscape(state),
			"Invalid credentials",
        )
        ctx.Redirect(http.StatusFound, loginURL)
        return
    }


    h.CompleteAuthorization(ctx, clientID, user.ID, redirectURI, scope, state)
}

func (h *authController) CompleteAuthorization(ctx *gin.Context, clientID, userID, redirectURI, scope, state string) {

    code,err := h.authSvc.StoreAuthCode(clientID, userID)
    if err != nil {
		cErr := cerror.NewError(http.StatusInternalServerError, err)
		response.ErrorResult(ctx,cErr)
        return
    }

    callback := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)
    ctx.Redirect(http.StatusFound, callback)
}
// func (h *authController) handleLogout(ctx *gin.Context) {
// 	userID, isExist := ctx.Get("AuthUserId")
// 	if !isExist {
// 		cErr := cerror.NewError(http.StatusUnauthorized, errors.New("Unauthorized")) 
// 		response.ErrorResult(ctx,cErr)
// 		return
// 	}

// 	err := h.authSvc.Logout(userID.(string))
// 	if err != nil {
// 		cErr := cerror.NewError(http.StatusInternalServerError, err) 
// 		response.ErrorResult(ctx,cErr)
// 		return
// 	}
// 	response.SuccessResult(ctx, http.StatusOK, nil)
// }

