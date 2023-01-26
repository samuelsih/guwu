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
	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/samuelsih/guwu/pkg/mail"
	"github.com/samuelsih/guwu/pkg/redis"
	"github.com/samuelsih/guwu/pkg/securer"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB *sqlx.DB
)

func TestMain(m *testing.M) {
	cleanup, err := setup()
	securer.SetSecret("0f5297b6f0114171e9de547801b1e8bb929fe1d091e63c6377a392ec1baa3d0b")

	if err != nil {
		log.Fatal(err)
	}

	if testDB == nil {
		log.Fatal("testDB is nil")
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
		input := RegisterInput{
			Username: "SevenSins",
			Password: "Seven123!",
		}

		expected := RegisterOutput{CommonResponse: business.CommonResponse{
			StatusCode: 400,
			Msg:        errEmailRequired.Error(),
		}}

		got := deps.Register(context.Background(), input, business.CommonInput{})

		if got != expected {
			t.Fatalf("TestRegister.RegisterEmptyEmail - expected %v, got %v", expected, got)
		}
	})

	t.Run("RegisterEmptyUsername", func(t *testing.T) {
		input := RegisterInput{
			Email:    "a@gmail.com",
			Password: "123123123",
		}

		expected := RegisterOutput{CommonResponse: business.CommonResponse{
			StatusCode: 400,
			Msg:        errUsernameRequired.Error(),
		}}

		got := deps.Register(context.Background(), input, business.CommonInput{})

		if got != expected {
			t.Fatalf("TestRegister.RegisterEmptyUsername - expected %v, got %v", expected, got)
		}
	})

	t.Run("RegisterEmptyPassword", func(t *testing.T) {
		input := RegisterInput{
			Email:    "a@gmail.com",
			Username: "123123123",
		}

		expected := RegisterOutput{CommonResponse: business.CommonResponse{
			StatusCode: 400,
			Msg:        errPasswordRequired.Error(),
		}}

		got := deps.Register(context.Background(), input, business.CommonInput{})

		if got != expected {
			t.Fatalf("TestRegister.RegisterEmptyUsername - expected %v, got %v", expected, got)
		}
	})

	t.Run("RegisterMultipleAcc", func(t *testing.T) {
		d := Deps{
			DB: testDB,
			SendEmail: func(ctx context.Context, param mail.Param, data any) error {
				return nil
			},
		}

		in := d.Register(context.Background(), RegisterInput{
			Email:    "testing@gmail.com",
			Username: "heavenlybrush",
			Password: "Heaven123!",
		},
			business.CommonInput{})

		if in.StatusCode != 200 {
			t.Fatalf("TestRegister.RegisterMultipleAcc - expected 200, got %v", in)
		}

		input := RegisterInput{
			Email:    "testing@gmail.com",
			Username: "heavenlybrush",
			Password: "Heaven123!",
		}

		expected := RegisterOutput{CommonResponse: business.CommonResponse{
			StatusCode: 400,
		}}

		got := d.Register(context.Background(), input, business.CommonInput{})

		if got.StatusCode != expected.StatusCode || !strings.Contains(got.Msg, "already taken") {
			t.Fatalf("TestRegister.RegisterEmptyUsername - expected %v, got %v", expected, got)
		}
	})
}

func TestLogin(t *testing.T) {
	deps := Deps{
		DB: testDB,
	}

	t.Parallel()

	t.Run("EmptyEmail", func(t *testing.T) {
		input := LoginInput{
			Password: "123123123",
		}

		expected := LoginOutput{
			CommonResponse: business.CommonResponse{
				StatusCode: 400,
				Msg:        errEmailRequired.Error(),
			},
		}

		got := deps.Login(context.Background(), input, business.CommonInput{})

		if expected != got {
			t.Fatalf("TestLogin.EmptyEmail - expected %v, got %v", expected, got)
		}
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		input := LoginInput{
			Email: "mail@gmail.com",
		}

		expected := LoginOutput{
			CommonResponse: business.CommonResponse{
				StatusCode: 400,
				Msg:        errPasswordRequired.Error(),
			},
		}

		got := deps.Login(context.Background(), input, business.CommonInput{})

		if expected != got {
			t.Fatalf("TestLogin.EmptyPassword - expected %v, got %v", expected, got)
		}
	})

	t.Run("UnknownUser", func(t *testing.T) {
		input := LoginInput{
			Email:    "ehehey@gmail.com",
			Password: "Mail123123!",
		}

		expected := LoginOutput{
			CommonResponse: business.CommonResponse{
				StatusCode: 400,
				Msg:        "unknown user",
			},
		}

		got := deps.Login(context.Background(), input, business.CommonInput{})

		if expected != got {
			t.Fatalf("TestLogin.UnknownUser - expected %v, got %v", expected, got)
		}
	})

	t.Run("Success", func(t *testing.T) {
		successDeps := Deps{
			DB: testDB,
			Store: func(ctx context.Context, key string, in any, time int64) error {
				return nil
			},
		}

		in := successDeps.Register(context.Background(), RegisterInput{
			Username: "gustalagusta",
			Email:    "gustalagusta@gmail.com",
			Password: "Testing123!",
		},

			business.CommonInput{})

		if in.StatusCode != 200 {
			t.Log("status code is not 200")
			t.Fatalf("TestLogin.Success - in should be 200, got %v", in)
		}

		input := LoginInput{
			Email:    "gustalagusta@gmail.com",
			Password: "Testing123!",
		}

		expected := LoginOutput{
			CommonResponse: business.CommonResponse{
				StatusCode: 200,
				Msg:        "OK",
			},
		}

		got := successDeps.Login(context.Background(), input, business.CommonInput{})

		if expected.StatusCode != got.StatusCode || expected.Msg != got.Msg || got.User == (model.User{}) {
			t.Fatalf("TestLogin.Success - expected %v, got %v", expected, got)
		}
	})

	t.Run("SuccessButErrorOnSession", func(t *testing.T) {
		successDeps := Deps{
			DB: testDB,
			Store: func(ctx context.Context, key string, in any, time int64) error {
				return errs.E(errs.Op("some_op"), errs.KindUnexpected, errors.New("error creating session"), "internal error")
			},
		}

		in := successDeps.Register(context.Background(), RegisterInput{
			Username: "andremaniani",
			Email:    "andre@gmail.com",
			Password: "Andre123!",
		},
			business.CommonInput{})

		if in.StatusCode != 200 {
			t.Fatalf("TestLogin.Success - in should be 200, got %v", in)
		}

		input := LoginInput{
			Email:    "andre@gmail.com",
			Password: "Andre123!",
		}

		expected := LoginOutput{
			CommonResponse: business.CommonResponse{
				StatusCode: 500,
				Msg:        redis.ErrInternal.Error(),
			},
		}

		got := successDeps.Login(context.Background(), input, business.CommonInput{})

		if expected.StatusCode != got.StatusCode || expected.Msg != got.Msg {
			t.Fatalf("TestLogin.SuccessButErrorOnSession - expected %v, got %v", expected, got)
		}
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()

	sessionEncrypted, err := securer.Encrypt([]byte("i-am-session"))
	if err != nil {
		t.Fatalf("TestLogout.Encrypt, got err: %v", err)
	}

	t.Run("EmptySessionID", func(t *testing.T) {
		deps := Deps{}

		out := deps.Logout(context.Background(), business.CommonInput{})

		if out.StatusCode != 400 {
			t.Fatalf("TestLogout.EmptySession - expected 400 got %d - %v", out.StatusCode, out)
		}
	})

	t.Run("UnknownSessionID", func(t *testing.T) {
		deps := Deps{
			Destroy: func(ctx context.Context, sessionID string) error {
				return errs.E(errs.Op("some_op"), errs.KindBadRequest, err, "unknown input")
			},
		}

		input := business.CommonInput{SessionID: sessionEncrypted}

		out := deps.Logout(context.Background(), input)

		if out.StatusCode != 400 {
			t.Fatalf("TestLogout.EmptySession - expected 400 got %d - %v", out.StatusCode, out)
		}
	})

	t.Run("InternalErr", func(t *testing.T) {
		internalErrDeps := Deps{
			Destroy: func(ctx context.Context, sessionID string) error {
				return errs.E(errs.Op("some_op"), errs.KindUnexpected, err, "unknown input")
			},
		}

		input := business.CommonInput{SessionID: sessionEncrypted}

		out := internalErrDeps.Logout(context.Background(), input)

		if out.StatusCode != 500 {
			t.Fatalf("TestLogout.InternalErr - expected 500 got %d - %v", out.StatusCode, out)
		}
	})

	t.Run("Success", func(t *testing.T) {
		deps := Deps{
			Destroy: func(ctx context.Context, sessionID string) error {
				return nil
			},
		}

		input := business.CommonInput{SessionID: sessionEncrypted}

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

	cleanup := func() error {
		return container.Terminate(ctx)
	}

	return cleanup, nil
}
