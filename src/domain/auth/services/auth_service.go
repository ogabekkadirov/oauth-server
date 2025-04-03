package services

import (
	"errors"

	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/helpers"
	"gitlab.com/yammt/oauth-auth-service/src/Infrastructure/jwt"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/repositories"
)

const (
	GrantTypeClientCredentials = "client_credentials"
	GrantTypePassword          = "password"
	GrantTypeAuthorizationCode = "authorization_code"	
	GrantTypeRefreshToken      = "refresh_token"
)

type AuthService interface {
	HandleClientCredentials(clientID, clientSecret string) (*models.Token, error)
	HandleAuthorizationCodeGrant(code, clientID, redirectURI string,clientSecret string) (*models.Token, error)
	HandlePasswordGrant(username, password, clientID string,clientSecret string) (*models.Token, error)
	HandleRefreshToken(refreshToken, clientID string,clientSecret string) (*models.Token, error)
	GetUserByID(userID string) (*models.User, error)
}


type authSvcImpl struct {
	userRepo repositories.UserRepository
    clientRepo  repositories.ClientRepository
    tokenRepo repositories.TokenRepository
    jwtGen jwt.JwtService
	authCodeRepo repositories.AuthCodeRepository
}

func NewAuthService(userRepo repositories.UserRepository,
					clientRepo  repositories.ClientRepository,
					tokenRepo repositories.TokenRepository,
					jwtGen jwt.JwtService,
					authCodeRepo repositories.AuthCodeRepository,
					) AuthService {
	return &authSvcImpl{
        userRepo: userRepo,
        clientRepo: clientRepo,
        tokenRepo: tokenRepo,
        jwtGen: jwtGen,
		authCodeRepo:authCodeRepo,
	}
}

func (s *authSvcImpl) HandleClientCredentials(clientID, clientSecret string) (*models.Token, error) {
	client, err := s.clientRepo.ValidateClient(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	if err := helpers.ValidateClientGrant(client.GrantTypes, GrantTypeClientCredentials); err != nil {
		return nil, err
	}
	tokenStr, err := s.jwtGen.GenerateAccessToken("", client.ID, client.Scopes)
	if err != nil {
		return nil, err
	}
	token := &models.Token{
		AccessToken: tokenStr,
		TokenType: "Bearer",
	}
	s.tokenRepo.StoreAccessToken(token, "")
	return token, nil
}

func (s *authSvcImpl) HandlePasswordGrant(username, password, clientID string,clientSecret string) (*models.Token, error) {
	client, err := s.clientRepo.ValidateClient(clientID,clientSecret)
	if err != nil {
		return nil, err
	}
	if err := helpers.ValidateClientGrant(client.GrantTypes, GrantTypePassword); err != nil {
		return nil, err
	}

	user, err := s.userRepo.ValidateUser(username, password)
	if err != nil {
		return nil, err
	}

	tokenStr, err := s.jwtGen.GenerateAccessToken(user.ID, clientID, []string{"read"})
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtGen.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	token := &models.Token{
		AccessToken: tokenStr,
		RefreshToken: refreshToken,
		TokenType: "Bearer",
	}
	s.tokenRepo.StoreAccessToken(token, user.ID)
	s.tokenRepo.StoreRefreshToken(refreshToken, user.ID)
	return token, nil
}

func (s *authSvcImpl) HandleAuthorizationCodeGrant(code, clientID, redirectURI string,clientSecret string) (*models.Token, error) {

	client, err := s.clientRepo.ValidateClient(clientID,clientSecret)
	if err != nil {
		return nil, err
	}
	if err := helpers.ValidateClientGrant(client.GrantTypes, GrantTypeAuthorizationCode); err != nil {
		return nil, err
	}

	userID, err := s.authCodeRepo.Validate(code)
	if err != nil {
		return nil, err
	}
	_ = s.authCodeRepo.Delete(code) 

	tokenStr, err := s.jwtGen.GenerateAccessToken(userID, clientID, []string{"read"})
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtGen.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	token := &models.Token{
		AccessToken: tokenStr,
		RefreshToken: refreshToken,
		TokenType: "Bearer",
	}
	s.tokenRepo.StoreAccessToken(token, userID)
	s.tokenRepo.StoreRefreshToken(refreshToken, userID)
	return token, nil
}

func (s *authSvcImpl) HandleRefreshToken(refreshToken, clientID string,clientSecret string) (*models.Token, error) {
	client, err := s.clientRepo.ValidateClient(clientID,clientSecret)
	if err != nil {
		return nil, err
	}
	if err := helpers.ValidateClientGrant(client.GrantTypes, GrantTypeRefreshToken); err != nil {
		return nil, err
	}

	userID, err := s.tokenRepo.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	tokenStr, err := s.jwtGen.GenerateAccessToken(userID, clientID, []string{"read"})
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := s.jwtGen.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	token := &models.Token{
		AccessToken: tokenStr,
		RefreshToken: newRefreshToken,
		TokenType: "Bearer",
	}
	s.tokenRepo.StoreAccessToken(token, userID)
	s.tokenRepo.StoreRefreshToken(newRefreshToken, userID)
	return token, nil
}


func (s *authSvcImpl) GetUserByID(userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}