package main

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func getEventKeys() ([]string, error) {
	redisHost := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

	keys, err := rdb.Keys(rdb.Context(), "*").Result()
	if err != nil {
		log.Error("Failed to get redis keys", err)
		return nil, err
	}

	return keys, nil
}

func getIssueReport(key string) (string, error) {
	redisHost := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

	report, err := rdb.Get(rdb.Context(), key).Result()
	if err != nil {
		log.Error("Failed to get report from redis", err)
		return "", err
	}

	return report, nil
}
