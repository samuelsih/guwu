package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
)

func NewRedis(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
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
