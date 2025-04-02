package seeder

import "gitlab.com/yammt/oauth-auth-service/src/Infrastructure/db"

func RunSeeder() {
	pool := db.NewPostgresPool()
	defer pool.Close()

	_ = SeedUsers(pool)
	_ = SeedClients(pool)
}