package mail

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var client Client

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatal("err client: ", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestSend(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err := client.Send(ctx, "info@gmail.com", "foo@gmail.com", "Hello", "World")
	if err != nil {
		e := err.(*errs.Error)
		t.Fatalf("success err is not nil: %v", e.Err)
	}
}

func TestSendMany(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	users := []struct {
		subject string
		body string
		to string
	}{
		{"Toni", "Tester", "toni.tester@example.com"},
		{"Tina", "Tester", "tina.tester@example.com"},
		{"John", "Doe", "john.doe@example.com"},
	}

	for _, user := range users {
		if err := client.Send(ctx, "info@gmail.com", user.to, user.subject, user.body); err != nil {
			e := err.(*errs.Error)
			t.Fatalf("success err is not nil: %v", e.Err)
		}
	}
}

func TestFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	t.Run("empty from", func(t *testing.T) {
		err := client.Send(ctx, "", "foo@gmail.com", "Hello", "World")
		if err == nil {
			t.Fatal("send from must fail")
		}
	})

	t.Run("empty to", func(t *testing.T) {
		err := client.Send(ctx, "info@gmail.com", "", "Hello", "World")
		if err == nil {
			t.Fatal("send to must fail")
		}
	})
}

func TestContextCancellation(t *testing.T) {
	c, cancelFunc := context.WithTimeout(context.Background(), time.Microsecond * 1)
	defer cancelFunc()

	err := client.Send(c, "info@gmail.com", "foo@gmail.com", "Hello", "World")
	if err == nil {
		t.Fatalf("context cancellation not hit, error is %v", err)
	}
}

func setup() error {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mailhog/mailhog:latest",
		ExposedPorts: []string{"1025/tcp"},
		WaitingFor:   wait.ForListeningPort("1025/tcp"),
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

	mappedPort, err := container.MappedPort(ctx, "1025")
	if err != nil {
		return err
	}
	port, _ := strconv.Atoi(mappedPort.Port())

	hostIP, err := container.Host(ctx)
	if err != nil {
		return err
	}

	client, err = NewClient(hostIP, port, "info@gmail.com", "", 20 * time.Second)
	return err
}
