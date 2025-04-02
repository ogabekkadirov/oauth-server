package redis

import (
	"context"
	"log"

	redis "github.com/redis/go-redis/v9"
)


type RedisClient struct {
	RdbClient *redis.Client
}



func NewRedisClient(Addr string) *RedisClient {
	rdbClient,err := NewClient(Addr)
		if err != nil {
		log.Println("ERROR: cannot connect to redis")
		panic(err)
	}
	return &RedisClient{
		RdbClient: rdbClient,
	}
}

func NewClient(Addr string) (rdb *redis.Client, err error) {
    rdb = redis.NewClient(&redis.Options{
        Addr:     Addr, // Redis server address
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	// Ping the Redis server to check if it's reachable
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("❌ cannot connect to redis:", err)
	}
	log.Println("✅ Connected to Postgres via pgx")
	return
}