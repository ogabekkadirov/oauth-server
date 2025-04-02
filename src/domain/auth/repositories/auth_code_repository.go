package repositories

import (
	"time"
)

type AuthCodeRepository interface {
	Save(code, userID string, ttl time.Duration) error
	Validate(code string) (string, error)
	Delete(code string) error
}