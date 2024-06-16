package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func RedisConnect() *redis.Client {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	REDIS_HOST := os.Getenv("REDIS_HOST")

	addr, err := redis.ParseURL(REDIS_HOST)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(addr)

	// opt, err := redis.ParseURL(REDIS_HOST)
	// if err != nil {
	// 	panic(err)
	// }

	// client := redis.NewClient(opt)
	return rdb
}
