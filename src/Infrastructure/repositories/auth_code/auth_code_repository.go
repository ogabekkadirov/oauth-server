package auth_code

import (
	"context"
	"errors"
	"time"

	"github.com/ogabekkadirov/oauth-server/src/domain/auth/repositories"
	redis "github.com/redis/go-redis/v9"
)

type authCodeRepoImpl struct {
	rdb *redis.Client
}

func NewAuthRepository(rdb *redis.Client) repositories.AuthCodeRepository{
	return &authCodeRepoImpl{rdb: rdb}
}	

func (r *authCodeRepoImpl) Save(code, userID string, ttl time.Duration) error {
	ctx := context.Background()
	return r.rdb.Set(ctx, "authcode:"+code, userID, ttl).Err()
}

func (r *authCodeRepoImpl) Validate(code string) (string, error) {
	ctx := context.Background()
	userID, err := r.rdb.Get(ctx, "authcode:"+code).Result()
	if err == redis.Nil {
		return "", errors.New("invalid or expired code")
	}
	return userID, err
}

func (r *authCodeRepoImpl) Delete(code string) error {
	ctx := context.Background()
	return r.rdb.Del(ctx, "authcode:"+code).Err()
}