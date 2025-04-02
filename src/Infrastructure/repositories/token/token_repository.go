package repositories

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/repositories"
)

const (
// userTable = "token"
)

type tokenRepoImpl struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) repositories.TokenRepository {
	return &tokenRepoImpl{rdb: rdb}
}

func (s *tokenRepoImpl) StoreAccessToken(token *models.Token, userID string) error {
	ctx := context.Background()
	return s.rdb.Set(ctx, "access:"+token.AccessToken, userID, time.Hour).Err()
}

func (s *tokenRepoImpl) StoreRefreshToken(refreshToken, userID string) error {
	ctx := context.Background()
	return s.rdb.Set(ctx, "refresh:"+refreshToken, userID, time.Hour*24*7).Err()
}

func (s *tokenRepoImpl) ValidateRefreshToken(refreshToken string) (string, error) {
	ctx := context.Background()
	return s.rdb.Get(ctx, "refresh:"+refreshToken).Result()
}