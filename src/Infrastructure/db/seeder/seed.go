package seeder

import "github.com/ogabekkadirov/oauth-server/src/Infrastructure/db"

func RunSeeder() {
	pool := db.NewPostgresPool()
	defer pool.Close()

	_ = SeedUsers(pool)
	_ = SeedClients(pool)
}