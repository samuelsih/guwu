package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB        *sqlx.DB
	testSessionDB *redis.Client
)

func TestMain(m *testing.M) {
	if err := setupDB(); err != nil {
		log.Fatal(err)
	}

	if err := setupRedis(); err != nil {
		log.Fatal(err)
	}

	if testDB == nil {
		log.Fatal("testdb is nil")
	}

	if testSessionDB == nil {
		log.Fatal("test session db is nil")
	}

	code := m.Run()
	os.Exit(code)
}

func setupDB() error {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
        ExposedPorts: []string{"5432/tcp"},
        WaitingFor:   wait.ForListeningPort("5432/tcp"),
        Env: map[string]string{
            "POSTGRES_DB":       "testdb",
            "POSTGRES_PASSWORD": "postgres",
            "POSTGRES_USER":     "postgres",
        },
	}

	container, err := testcontainers.GenericContainer(
			ctx, 
			testcontainers.GenericContainerRequest{
        	ContainerRequest: req,
        	Started:          true,
		},
	)

	if err != nil {
		return err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
    if err != nil {
        return err
    }

	hostIP, err := container.Host(ctx)
    if err != nil {
        return err
    }

	uri := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", hostIP, mappedPort.Port())

	testDB = config.ConnectPostgres(uri)
	if testDB == nil {
		return errors.New("cannot connect testGuestDB")
	}

	if err := config.MigrateAll(testDB); err != nil {
		return err
	}
	
	return nil
}

func setupRedis() error {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())
	testSessionDB = config.NewRedis(uri)

	return nil
}