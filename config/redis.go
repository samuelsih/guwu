package config

import (
	"context"
	"log"

	"github.com/rueian/rueidis"
)

func NewRedis(host, password string) rueidis.Client {
	conn, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{host},
		Password:    password,
	})

	if err != nil {
		log.Printf("cant create redis connection: %v", err)
		return nil
	}

	err = conn.Do(context.Background(), conn.B().Ping().Build()).Error()
	if err != nil {
		log.Printf("cant ping redis: %v", err)
		return nil
	}

	log.Println("Redis connect")

	return conn
}
