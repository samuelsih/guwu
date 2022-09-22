package config

import (
	"context"

	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
)

func NewRedis(url string) *redis.Client {
	var rdb *redis.Client
	
	if url == "" {
		url = `127.0.0.1:6379`
		rdb = redis.NewClient(&redis.Options{
			Addr: url,
			Password: "",
			DB: 0,
		})
	} else {
		opt, err := redis.ParseURL(url)
		if err != nil {
			panic(err)
		}
	
		rdb = redis.NewClient(opt)
	}

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return rdb
}
