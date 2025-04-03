package repositories

import "github.com/ogabekkadirov/oauth-server/src/domain/auth/models"

type ClientRepository interface {
	ValidateClient(id, secret string) (*models.Client, error)
	GetByID(id string) (*models.Client, error)
}