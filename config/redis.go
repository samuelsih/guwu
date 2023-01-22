package config

import (
	"context"

	"github.com/rueian/rueidis"
	"github.com/samuelsih/guwu/pkg/logger"
)

func NewRedis(host, password string) rueidis.Client {
	conn, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{host},
		Password:    password,
	})

	if err != nil {
		logger.SysFatal("cant create redis connection: %v", err)
		return nil
	}

	err = conn.Do(context.Background(), conn.B().Ping().Build()).Error()
	if err != nil {
		logger.SysFatal("cant ping redis: %v", err)
		return nil
	}

	logger.SysInfo("Redis connect")

	return conn
}
