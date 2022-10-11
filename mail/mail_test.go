package mail

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var serverMailTest *Mailer

func TestMain(m *testing.M) {
	if err := setupMailServer(); err != nil {
		log.Fatal(err)
	}
}

func TestMailSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	defer serverMailTest.Close()

	msg := Message{
		To: "something@gmail.com",
		Subject: "HAI",
		PlainContent: "HAI",
		HTMLContent: "<h1>HAI</h1",
	}

	serverMailTest.Send(ctx, msg)

	select {
		case <-ctx.Done():
			t.Error("context done")

		case err := <- serverMailTest.errChan:
			if err != nil { t.Error(err) }
	}
}

func TestMailConcurrent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	defer serverMailTest.Close()

	msg := Message{
		To: "something@gmail.com",
		Subject: "HAI",
		PlainContent: "HAI",
		HTMLContent: "<h1>HAI</h1",
	}

	for i := 10; i < 10; i++ {
		go func() { serverMailTest.Send(ctx, msg) }()
	}

	select {
	case <-ctx.Done():
		t.Error("context done")

	case err := <- serverMailTest.errChan:
		if err != nil { t.Error(err) }
	}
}

func TestMailConcurrentFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	defer serverMailTest.Close()

	msg := Message{
		To: "something@gmail.com",
		Subject: "HAI",
		PlainContent: "HAI",
		HTMLContent: "<h1>HAI</h1",
	}

	for i := 10; i < 10; i++ {
		go func() { serverMailTest.Send(ctx, msg) }()
	}
	
	if ctx.Err() == nil {
		t.Error(`ctx should exit the send`)
	}
}


func setupMailServer() error {
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

	hostIP, err := container.Host(ctx)
	if err != nil {
		return err
	}

	host := fmt.Sprintf("mailhog://%s%s", hostIP, mappedPort.Port())
	serverMailTest = NewMailer(host, "localhost", "info@gmail.com", "")

	if serverMailTest == nil {
		return errors.New(`server is nil`)
	}

	return nil
}