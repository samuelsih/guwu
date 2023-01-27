package follow

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/samuelsih/guwu/pkg/securer"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sqlx.DB

var users = [2]struct{
	Username string
	Email string
	Password string
} {
	{
		Username: "adelombok",
		Email: "adelombok@gmail.com",
		Password: "$2a$07$GsdzeF04uKNmPyEf1R.WUOZF.i9Xhpx6peu3NBMN7NdPe//tWEfY",
	},

	{
		Username: "budijakarta",
		Email: "budijakarta@gmail.com",
		Password: "$2a$07$xrzfGJNLhcl/nzedzARjKOQLHuFZRmd0hn6Z/Mdlj2hzxkIgyfLfm",
	},
}

func TestMain(m *testing.M) {
	cleanup, err := setup()
	securer.SetSecret("0f5297b6f0114171e9de547801b1e8bb929fe1d091e63c6377a392ec1baa3d0b")

	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := cleanup(); err != nil {
		log.Fatalf("error cleaning up: %v", err)
	}

	os.Exit(code)
}

func TestFollow(t *testing.T) {
	t.Parallel()

	t.Run("Unauthenticated", func(t *testing.T) {
		deps := Deps {
			DB: testDB,
		}

		out := deps.Follow(context.Background(), FollowIn{}, business.CommonInput{})

		if out.StatusCode != 403 {
			t.Fatalf("expected status code 400, got %d - %v", out.StatusCode, out)
		}
	})

	t.Run("Unknown session id", func(t *testing.T) {
		sess, _ := securer.Encrypt([]byte("1231231231231232123"))

		deps := Deps {
			DB: testDB,
			GetUserSession: func(ctx context.Context, key string, dst any) error {
				return errs.E(errs.Op("GetUserSession"), errs.KindBadRequest, errors.New("unknown input"), "unknown input")
			},
		}

		in := FollowIn{
			UserID: "123123123",
		}

		cmn := business.CommonInput{
			SessionID: sess,
		}

		out := deps.Follow(context.Background(), in, cmn)
		expected := FollowOut {
			business.CommonResponse{
				StatusCode: 400,
				Msg: "unknown input",
			},
		}

		if out != expected {
			t.Fatalf("expected %v, got %v", expected, out)
		}
	})

	t.Run("Unknown user_follow_id", func(t *testing.T) {
		sess, _ := securer.Encrypt([]byte("1231231231231232123"))
		deps := Deps {
			DB: testDB,
			GetUserSession: func(ctx context.Context, key string, dst any) error {
				return nil
			},
		}

		in := FollowIn{
			UserID: "123123123",
		}

		cmn := business.CommonInput{
			SessionID: sess,
		}

		out := deps.Follow(context.Background(), in, cmn)
		expected := FollowOut {
			business.CommonResponse{
				StatusCode: 400,
				Msg: "unknown id for follows user id",
			},
		}

		if out != expected {
			t.Fatalf("expected %v, got %v", expected, out)
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

	qs := `INSERT INTO users(username, email, password)
	VALUES (:username, :email, :password)
	`

	r, err := testDB.NamedExecContext(ctx, qs, users)
	if err != nil {
		return nil, err
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected != 2 {
		return nil, errors.New("affected not 2, but " + fmt.Sprint(affected))
	}

	cleanup := func() error {
		return container.Terminate(ctx)
	}

	return cleanup, nil
}
