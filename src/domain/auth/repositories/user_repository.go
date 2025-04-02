package repositories

import "gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	ValidateUser(username, password string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Create(user *models.User) error
}