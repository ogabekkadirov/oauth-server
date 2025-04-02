package repositories

import "gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"

type ClientRepository interface {
	ValidateClient(id, secret string) (*models.Client, error)
	GetByID(id string) (*models.Client, error)
}