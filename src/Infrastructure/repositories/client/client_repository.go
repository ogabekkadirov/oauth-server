package client

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/models"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/repositories"
)

type clientRepoImpl struct {
    db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) repositories.ClientRepository {
    return &clientRepoImpl{db: db}
}

func (r *clientRepoImpl) ValidateClient(id, secret string) (*models.Client, error) {
	ctx := context.Background()
	var client models.Client
	err := r.db.QueryRow(ctx, `SELECT id, secret,redirect_uris,grant_types FROM clients WHERE id=$1 AND secret=$2`, id, secret).
		Scan(&client.ID, &client.Secret,&client.RedirectURIs, &client.GrantTypes)
	if err != nil {
		return nil, errors.New("invalid client credentials")
	}
	client.Scopes = []string{"read"} // simplify
	return &client, nil
}

func (r *clientRepoImpl) GetByID(clientID string) (*models.Client, error) {
    ctx := context.Background()
	var client models.Client
	err := r.db.QueryRow(ctx, `SELECT id,secret,redirect_uris,grant_types FROM clients WHERE id=$1`, clientID).
		Scan(&client.ID, &client.Secret, &client.RedirectURIs, &client.GrantTypes)
	if err != nil {
		return nil, err
	}
	client.Scopes = []string{"read"}
	return &client, nil
}

