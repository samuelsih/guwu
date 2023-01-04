package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/model"
	"github.com/samuelsih/guwu/pkg/session"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB *sqlx.DB
)

func TestMain(m *testing.M) {
	cleanup, err := setup()

	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := cleanup(); err != nil {
		log.Fatalf("could not cleanup : %v", err)
	}

	os.Exit(code)
}


func TestRegister(t *testing.T) {
	deps := Deps {
		DB: testDB,
	}

	t.Parallel()

	t.Run("RegisterEmptyEmail", func(t *testing.T) {
		input := RegisterInput {
			Username: "SevenSins",
			Password: "Seven123!",
		}

		expected := RegisterOutput {CommonResponse: business.CommonResponse{
			StatusCode: 400,
			Msg: errEmailRequired.Error(),
		}}

		got := deps.Register(context.Background(), input)

		if got != expected {
			t.Fatalf("TestRegister.RegisterEmptyEmail - expected %v, got %v", expected, got)
		}
	})

	t.Run("RegisterEmptyUsername", func(t *testing.T) {
		input := RegisterInput {
			Email: "a@gmail.com",
			Password: "123123123",
		}

		expected := RegisterOutput {CommonResponse: business.CommonResponse{
			StatusCode: 400,
			Msg: errUsernameRequired.Error(),
		}}

		got := deps.Register(context.Background(), input)

		if got != expected {
			t.Fatalf("TestRegister.RegisterEmptyUsername - expected %v, got %v", expected, got)
		}
	})

	t.Run("RegisterEmptyPassword", func(t *testing.T) {
		input := RegisterInput {
			Email: "a@gmail.com",
			Username: "123123123",
		}

		expected := RegisterOutput {CommonResponse: business.CommonResponse{
			StatusCode: 400,
			Msg: errPasswordRequired.Error(),
		}}

		got := deps.Register(context.Background(), input)

		if got != expected {
			t.Fatalf("TestRegister.RegisterEmptyUsername - expected %v, got %v", expected, got)
		}
	})


	t.Run("RegisterMultipleAcc", func(t *testing.T) {
		in := deps.Register(context.Background(), RegisterInput{
			Email: "testing@gmail.com",
			Username: "heavenlybrush",
			Password: "Heaven123!",
		})

		if in.StatusCode != 200 {
			t.Fatalf("TestRegister.RegisterMultipleAcc - expected 200, got %v", in)
		}

		input := RegisterInput {
			Email: "testing@gmail.com",
			Username: "heavenlybrush",
			Password: "Heaven123!",
		}

		expected := RegisterOutput {CommonResponse: business.CommonResponse{
			StatusCode: 400,
		}}

		got := deps.Register(context.Background(), input)

		if got.StatusCode != expected.StatusCode || !strings.Contains(got.Msg, "already taken") {
			t.Fatalf("TestRegister.RegisterEmptyUsername - expected %v, got %v", expected, got)
		}
	})
}

func TestLogin(t *testing.T) {
	deps := Deps {
		DB: testDB,
	}

	t.Parallel()

	t.Run("EmptyEmail", func(t *testing.T) {
		input := LoginInput {
			Password: "123123123",
		}

		expected := LoginOutput {
			CommonResponse: business.CommonResponse{
				StatusCode: 400,
				Msg: errEmailRequired.Error(),
			},
		}

		got := deps.Login(context.Background(), input)

		if expected != got {
			t.Fatalf("TestLogin.EmptyEmail - expected %v, got %v", expected, got)
		}
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		input := LoginInput {
			Email: "mail@gmail.com",
		}

		expected := LoginOutput {
			CommonResponse: business.CommonResponse{
				StatusCode: 400,
				Msg: errPasswordRequired.Error(),
			},
		}

		got := deps.Login(context.Background(), input)

		if expected != got {
			t.Fatalf("TestLogin.EmptyPassword - expected %v, got %v", expected, got)
		}
	})

	t.Run("UnknownUser", func(t *testing.T) {
		input := LoginInput {
			Email: "mail@gmail.com",
			Password: "Mail123123!",
		}

		expected := LoginOutput {
			CommonResponse: business.CommonResponse{
				StatusCode: 400,
				Msg: "unknown user",
			},
		}

		got := deps.Login(context.Background(), input)

		if expected != got {
			t.Fatalf("TestLogin.UnknownUser - expected %v, got %v", expected, got)
		}
	})

	t.Run("Success", func(t *testing.T) {
		successDeps := Deps {
			DB: testDB,
			CreateSession: func(ctx context.Context, in any) (string, error) {
				return "123", nil
			},
		}

		in := successDeps.Register(context.Background(), RegisterInput {
			Username: "gustalagusta",
			Email: "gustalagusta@gmail.com",
			Password: "Testing123!",
		})

		if in.StatusCode != 200 {
			t.Fatalf("TestLogin.Success - in should be 200, got %v", in)
		}

		input := LoginInput {
			Email: "gustalagusta@gmail.com",
			Password: "Testing123!",
		}

		expected := LoginOutput {
			CommonResponse: business.CommonResponse{
				StatusCode: 200,
				Msg: "OK",
				SessionID: "123",
			},
		}

		got := successDeps.Login(context.Background(), input)

		if expected.StatusCode != got.StatusCode || expected.Msg != got.Msg || got.User == (model.User{}) || expected.SessionID != "123" {
			t.Fatalf("TestLogin.Success - expected %v, got %v", expected, got)
		}
	})

	t.Run("SuccessButErrorOnSession", func(t *testing.T) {
		successDeps := Deps {
			DB: testDB,
			CreateSession: func(ctx context.Context, in any) (string, error) {
				return "", session.InternalErr
			},
		}

		in := successDeps.Register(context.Background(), RegisterInput {
			Username: "andremaniani",
			Email: "andre@gmail.com",
			Password: "Andre123!",
		})

		if in.StatusCode != 200 {
			t.Fatalf("TestLogin.Success - in should be 200, got %v", in)
		}

		input := LoginInput {
			Email: "andre@gmail.com",
			Password: "Andre123!",
		}

		expected := LoginOutput {
			CommonResponse: business.CommonResponse{
				StatusCode: 500,
				Msg: session.InternalErr.Error(),
				SessionID: "",
			},
		}

		got := successDeps.Login(context.Background(), input)

		if expected != got {
			t.Fatalf("TestLogin.SuccessButErrorOnSession - expected %v, got %v", expected, got)
		}
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()

	t.Run("EmptySessionID", func(t *testing.T) {
		deps := Deps{}

		out := deps.Logout(context.Background(), LogoutInput{})

		if out.StatusCode != 400 {
			t.Fatalf("TestLogout.EmptySession - expected 400 got %d - %v", out.StatusCode, out)
		}
	})

	t.Run("UnknownSessionID", func(t *testing.T) {
		deps := Deps {
			DestroySession: func(ctx context.Context, sessionID string) error {
				return session.UnknownSessionID
			},
		}

		input := LogoutInput{
			CommonRequest: business.CommonRequest{
				SessionID: "123",
			},
		}

		out := deps.Logout(context.Background(), input)

		if out.StatusCode != 400 {
			t.Fatalf("TestLogout.EmptySession - expected 400 got %d - %v", out.StatusCode, out)
		}
	})

	t.Run("InternalErr", func(t *testing.T) {
		deps := Deps {
			DestroySession: func(ctx context.Context, sessionID string) error {
				return session.InternalErr
			},
		}

		input := LogoutInput{
			CommonRequest: business.CommonRequest{
				SessionID: "123",
			},
		}

		out := deps.Logout(context.Background(), input)

		if out.StatusCode != 500 {
			t.Fatalf("TestLogout.InternalErr - expected 500 got %d - %v", out.StatusCode, out)
		}
	})

	t.Run("Success", func(t *testing.T) {
		deps := Deps {
			DestroySession: func(ctx context.Context, sessionID string) error {
				return nil
			},
		}

		input := LogoutInput{
			CommonRequest: business.CommonRequest{
				SessionID: "02917joaisdd8v92b3ir",
			},
		}

		out := deps.Logout(context.Background(), input)

		if out.StatusCode != 200 {
			t.Fatalf("TestLogout.Success - expected 200 got %d - %v", out.StatusCode, out)
		}
	})
}

func setup() (func() error, error) {
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
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
    if err != nil {
        return nil, err
    }

	hostIP, err := container.Host(ctx)
    if err != nil {
        return nil, err
    }

	uri := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb?sslmode=disable", hostIP, mappedPort.Port())

	testDB = config.ConnectPostgres(uri)
	if testDB == nil {
		return nil, errors.New("cannot connect testGuestDB")
	}

	if err := config.LoadPostgresExtension(testDB); err != nil {
		return nil, errors.New("cannot load postgres extension")
	}

	if err := config.MigrateAll(testDB); err != nil {
		return nil, err
	}
	
	cleanup := func () error {
		return container.Terminate(ctx)
	}

	return cleanup, nil
}