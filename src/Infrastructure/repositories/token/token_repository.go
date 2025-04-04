package repositories

import (
	"context"
	"time"

	"github.com/ogabekkadirov/oauth-server/src/domain/auth/models"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/repositories"
	redis "github.com/redis/go-redis/v9"
)

const (
	ttl_access = time.Hour
	ttl_refresh = time.Hour * 24 * 7
)

type tokenRepoImpl struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) repositories.TokenRepository {
	return &tokenRepoImpl{rdb: rdb}
}

func (s *tokenRepoImpl) StoreAccessToken(token *models.Token, userID string) error {
	ctx := context.Background()
	return s.rdb.Set(ctx, "access:"+token.AccessToken, userID, ttl_access).Err()
}

func (s *tokenRepoImpl) StoreRefreshToken(refreshToken, userID string) error {
	ctx := context.Background()
	return s.rdb.Set(ctx, "refresh:"+refreshToken, userID, ttl_refresh).Err()
}

func (s *tokenRepoImpl) ValidateRefreshToken(refreshToken string) (string, error) {
	ctx := context.Background()
	return s.rdb.Get(ctx, "refresh:"+refreshToken).Result()
}