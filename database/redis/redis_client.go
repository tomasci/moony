package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	// to ensure instance created only once
	once sync.Once
	// redis client instance
	rdb *redis.Client
	// for errors
	initError error
)

func InitializeRedis() {
	once.Do(func() {
		log.Println("Initializing Redis client")

		// get keys from env
		host := os.Getenv("REDIS_HOST")
		db := os.Getenv("REDIS_DB")
		user := os.Getenv("REDIS_USER")
		password := os.Getenv("REDIS_PASSWORD")

		var dbInt int
		dbInt, initError = strconv.Atoi(db)
		if initError != nil {
			log.Fatalf("error parsing REDIS_DB: %v", initError)
		}

		rdb = redis.NewClient(&redis.Options{
			Addr:     host,
			Username: user,
			Password: password,
			DB:       dbInt,
			Protocol: 3,
		})

		_, initError = rdb.Ping(context.Background()).Result()
		if initError != nil {
			log.Fatalf("failed to connect to redis: %v", initError)
		}
	})
}

func GetRedisClient() (*redis.Client, error) {
	if rdb == nil {
		InitializeRedis()
	}

	return rdb, initError
}
