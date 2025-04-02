package seeder

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"
)

func SeedClients(db *pgxpool.Pool) error {
    ctx := context.Background()
    var exists bool

    err := db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM clients WHERE id = $1)`, "service-a").Scan(&exists)
    if err != nil {
        return err
    }

    if exists {
        fmt.Println("🟡 Client 'service-a' already exists, skipping seeding.")
        return nil
    }
    clients := []models.Client{
		{
			ID:           "client1",
			Secret:       "secret1",
			RedirectURIs: []string{"http://localhost:3000/callback"},
			GrantTypes:   []string{"client_credentials"},
		},
		{
			ID:           "client2",
			Secret:       "secret2",
			RedirectURIs: []string{"http://localhost:3000/callback"},
			GrantTypes:   []string{"password"},
		},
		{
			ID:           "client3",
			Secret:       "secret3",
			RedirectURIs: []string{"http://localhost:3000/callback"},
			GrantTypes:   []string{"authorization_code"},
		},
		{
			ID:           "client4",
			Secret:       "secret4",
			RedirectURIs: []string{"http://localhost:3000/callback"},
			GrantTypes:   []string{"refresh_token"},
		},
		{
			ID:           "client5",
			Secret:       "secret5",
			RedirectURIs: []string{"http://localhost:3000/callback"},
			GrantTypes:   []string{"client_credentials", "authorization_code", "refresh_token"},
		},
	}
    
	for _, client := range clients {
		var exists bool
		err := db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM clients WHERE id = $1)`, client.ID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("checking client existence failed: %w", err)
		}

		if exists {
			fmt.Printf("🟡 Client '%s' already exists, skipping.\n", client.ID)
			continue
		}

		_, err = db.Exec(ctx, `
			INSERT INTO clients (id, secret, redirect_uris, grant_types)
			VALUES ($1, $2, $3, $4)
		`, client.ID, client.Secret, client.RedirectURIs, client.GrantTypes)

		if err != nil {
			return fmt.Errorf("inserting client '%s' failed: %w", client.ID, err)
		}

		fmt.Printf("✅ Client '%s' seeded.\n", client.ID)
	}

    return nil
}
