package repositories

import "github.com/ogabekkadirov/oauth-server/src/domain/auth/models"

type TokenRepository interface {
	StoreAccessToken(token *models.Token, userID string) error
	StoreRefreshToken(refreshToken, userID string) error
	ValidateRefreshToken(refreshToken string) (string, error)
}
