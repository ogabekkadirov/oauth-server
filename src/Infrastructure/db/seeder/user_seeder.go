package seeder

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedUsers(db *pgxpool.Pool) error {
    ctx := context.Background()
    var exists bool

    err := db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, "admin").Scan(&exists)
    if err != nil {
        return err
    }

    if exists {
        fmt.Println("🟡 User 'admin' already exists, skipping seeding.")
        return nil
    }

    _, err = db.Exec(ctx, `
        INSERT INTO users (id, username, password)
        VALUES ($1, $2, $3)
    `, uuid.New().String(), "admin", "$2y$04$hNpTi/ynJT6k63LUkJR8hO0zx6g.EKLuafjWiyNhKf5z9sygJqF6y")

    if err != nil {
        return err
    }

    fmt.Println("✅ User 'admin' seeded.")
    return nil
}
