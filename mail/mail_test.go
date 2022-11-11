package mail

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

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

func TestMail(t *testing.T) {
	defer serverMailTest.Close()

	msg := Message{
		To: "something@gmail.com",
		Subject: "HAI",
		PlainContent: "HAI",
		HTMLContent: "<h1>HAI</h1",
	}

	serverMailTest.Send(msg)
}

func TestMailConcurrent(t *testing.T) {
	defer serverMailTest.Close()

	msg := Message{
		To: "something@gmail.com",
		Subject: "HAI",
		PlainContent: "HAI",
		HTMLContent: "<h1>HAI</h1",
	}

	for i := 0; i < 10; i++ {
		serverMailTest.Send(msg)
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

	port, _ := strconv.Atoi(mappedPort.Port())

	host := fmt.Sprintf("mailhog://%s%s", hostIP, mappedPort.Port())
	serverMailTest = NewMailer(host, port, "info@gmail.com", "")

	if serverMailTest == nil {
		return errors.New(`server is nil`)
	}

	return nil
}

