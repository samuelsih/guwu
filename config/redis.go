package config

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
)

func NewRedis() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("UPSTASH_URL"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Redis ready!")

	return rdb
}