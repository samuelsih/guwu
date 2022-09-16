package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/config"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testGuestDB *sqlx.DB
	testGuestSessionDB *redis.Client
)

func TestMain(m *testing.M) {
	if err := setupDB(); err != nil {
		log.Fatal(err)
	}

	if err := setupRedis(); err != nil {
		log.Fatal(err)
	}

	if testGuestDB == nil {
		log.Fatal("testdb is nil")
	}

	if testGuestSessionDB == nil {
		log.Fatal("test session db is nil")
	}

	code := m.Run()
	os.Exit(code)
}

func TestGuestRegister(t *testing.T) {
	ctx := context.Background()

	guest := Guest{
		DB: testGuestDB,
		SessionDB: testGuestSessionDB,
	}

	t.Run(`Register Empty struct`, func(t *testing.T) {
		in := GuestRegisterIn{}

		expected := GuestRegisterOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: errEmailRequired.Error(),
			},
		}

		out := guest.Register(ctx, &in)

		assert.Equal(t, expected, out)
		assert.Nil(t, out.User)
	})

	t.Run(`Insert Success`, func(t *testing.T) {
		in := GuestRegisterIn{
			Email: "samuel@gmail.com",
			Username: "SamuelSamuelSamuel",
			Password: "Akubohong123!",
		}

		expected := GuestRegisterOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		out := guest.Register(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
		assert.NotNil(t, out.User)
	})

	t.Run(`Duplicate Insert`, func(t *testing.T) {
		in := GuestRegisterIn{
			Email: "samuel@gmail.com",
			Username: "HayangUlinAing",
			Password: "GajadiBohong123!",
		}

		expected := GuestRegisterOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `email already exists`,
			},
		}

		out := guest.Register(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	guest := Guest{
		DB: testGuestDB,
		SessionDB: testGuestSessionDB,
	}

	t.Run(`Login Empty struct`, func(t *testing.T) {
		in := GuestLoginIn{}

		expected := GuestLoginOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `invalid email`,
			},
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
		assert.Nil(t, out.User)
	})

	t.Run(`Login Empty Password`, func(t *testing.T){
		in := GuestLoginIn{
			Email: "aiaiai@gmail.com",
		}

		expected := GuestLoginOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `invalid password`,
			},
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
		assert.Nil(t, out.User)
	})

	t.Run(`Login Email Not Found`, func(t *testing.T) {
		in := GuestLoginIn{
			Email: "bukansamuel@gmail.com",
			Password: "Passwordajah",
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, http.StatusBadRequest, out.StatusCode)
		assert.Nil(t, out.User)
	})

	t.Run(`Login With Wrong Password`, func(t *testing.T) {
		in := GuestLoginIn{
			Email: "samuel@gmail.com",
			Password: "123123",
		}

		expected := GuestLoginOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `email or password does not match`,
			},
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
		assert.Nil(t, out.User)
	})

	t.Run(`Success`, func(t *testing.T){
		in := GuestLoginIn{
			Email: "samuel@gmail.com",
			Password: "Akubohong123!",
		}

		expected := GuestLoginOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
		assert.NotNil(t, out.User)
		assert.NotEmpty(t, out.User.ID)
		assert.Equal(t, out.User.Email, "samuel@gmail.com")
		assert.Equal(t, out.User.Username, "SamuelSamuelSamuel")
	})
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

	testGuestDB = config.ConnectPostgres(uri)
	if testGuestDB == nil {
		return errors.New("cannot connect testGuestDB")
	}

	if err := config.MigrateAll(testGuestDB); err != nil {
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
	testGuestSessionDB = config.NewRedis(uri)

	return nil
}