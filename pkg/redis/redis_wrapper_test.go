package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/samuelsih/guwu/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var client Client

func TestMain(m *testing.M) {
	err := setup()

	if err != nil {
		log.Fatalf("setup: %v", err)
	}

	code := m.Run()

	client.Pool.Close()

	os.Exit(code)
}

func TestGetAndSet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	err := client.Set(ctx, "name", "budi", 100)
	if err != nil {
		t.Fatalf("Set: expected err is nil, got %v", err)
	}

	var result string
	result, err = client.Get(ctx, "name")
	if err != nil {
		t.Fatalf("Set: expected err is nil, got %v", err)
	}

	if result != "budi" {
		t.Fatalf("Result: expected budi, got %s", result)
	}
}

func TestGetAndSet_JSON(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type foo struct {
		Bar string `json:"bar"`
	}

	input := foo{
		Bar: "baz",
	}

	err := client.SetJSON(ctx, "struct", input, 100)
	if err != nil {
		t.Fatalf("Set: expected err is nil, got %v", err)
	}

	var result foo

	err = client.GetJSON(ctx, "struct", &result)
	if err != nil {
		t.Fatalf("Set: expected err is nil, got %v", err)
	}

	if result != input {
		t.Fatalf("Result: expected %v, got %s", input, result)
	}
}

func TestDestroy(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		err := client.Set(ctx, "age", "18", 100)
		if err != nil {
			t.Fatalf("Set: expected err is nil, got %v", err)
		}

		err = client.Destroy(ctx, "age")
		if err != nil {
			t.Fatalf("Set: expected err is nil, got %v", err)
		}
	})

	t.Run("unknown key", func(t *testing.T) {
		err := client.Set(ctx, "foosha", "ayooo", 100)
		if err != nil {
			t.Fatalf("Set: expected err is nil, got %v", err)
		}

		err = client.Destroy(ctx, "foozhaa")
		if err == nil || !errors.Is(err, ErrUnknownKey) {
			t.Fatalf("Destroy: expected err is ErrUnknownKey, got %v", err)
		}
	})
}

func setup() error {
	req := testcontainers.ContainerRequest{
		Image:        "redis",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

	ctx := context.Background()

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

	uri := fmt.Sprintf("%s:%s", hostIP, mappedPort.Port())

	db := config.NewRedis(uri, "")
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	client = *NewClient(db)

	return nil
}
