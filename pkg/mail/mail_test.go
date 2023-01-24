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

	p := Param {
		Name: "Foo",
		Email: "foo@gmail.com",
		Subject: "Hello",
		TemplateTypes: OTPMsg,
	}

	tplData := OTPTplData {
		Username: "Agus",
		OTP: "1234",
	}

	err := client.Send(ctx, p, tplData)
	if err != nil {
		e := err.(*errs.Error)
		t.Fatalf("success err is not nil: %v", e.Err)
	}
}

func TestSendMany(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	params := []Param {
		{
			Name: "Toni",	
			Email: "toni.tester@example.com",
			Subject: "OTP",
		},

		{
			Name: "Tina",	
			Email: "tina.tester@example.com",
			Subject: "OTP",
		},

		{
			Name: "John",	
			Email: "john.tester@example.com",
			Subject: "OTP",
		},
	}

	tplData := []any {
		OTPTplData{
			Username: "Toni",
			OTP: "1234",
		},

		OTPTplData {
			Username: "Tina",
			OTP: "4567",
		},

		OTPTplData {
			Username: "John",
			OTP: "0928",
		},
	}

	for i := 0; i < len(params); i++ {
		if err := client.Send(ctx, params[i], tplData[i]); err != nil {
			e := err.(*errs.Error)
			t.Fatalf("success err is not nil: %v", e.Err)	
		}
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	p := Param {
		Name: "Foo",
		Email: "foo@gmail.com",
		Subject: "Hello",
		TemplateTypes: OTPMsg,
	}

	tplData := OTPTplData {
		Username: "Agus",
		OTP: "1234",
	}

	err := client.Send(ctx, p, tplData)
	if err == nil {
		t.Fatal("context deadline not hit")
	}
}

func TestFailParam(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	t.Run("mismatch type", func(t *testing.T) {
		param := Param{
			Name: "foo",
			Email: "foo@gmail.com",
			Subject: "Fooo",
			TemplateTypes: OTPMsg,
		}

		tplData := RecoverPasswdTplData {
			Username: "bar",
			GeneratedLink: "localhost:something",
		}

		if err := client.Send(ctx, param, tplData); err == nil {
			t.Fatal("should not be nil")
		}
	})

	t.Run("empty param", func(t *testing.T) {
		tplData := RecoverPasswdTplData {
			Username: "bar",
			GeneratedLink: "localhost:something",
		}

		if err := client.Send(ctx, Param{}, tplData); err == nil {
			t.Fatal("empty param should result error")
		}
	})

	t.Run("unknown msg tpl type", func(t *testing.T) {
		param := Param{
			Name: "foo",
			Email: "foo@gmail.com",
			Subject: "Fooo",
			TemplateTypes: MsgType(10),
		}

		tplData := OTPTplData {
			Username: "Agus",
			OTP: "1234",
		}

		if err := client.Send(ctx, param, tplData); err == nil {
			t.Fatal("unknown type should result error")
		}
	})
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

	client, err = NewClient(hostIP, port, "info@gmail.com", "", "Guwu", "info@gmail.com")
	return err
}
