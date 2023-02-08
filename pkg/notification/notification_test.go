package notification

import (
	"log"
	"os"
	"testing"
)

var c Config

const (
	//https://github.com/pusher/push-notifications-go/blob/764224c311b854e5a272f8601b98957448a71995/push_notification_integration_test.go#L16
	instanceTest = "9aa32e04-a212-44ab-a592-9aeba66e46ac"

	//https://github.com/pusher/push-notifications-go/blob/764224c311b854e5a272f8601b98957448a71995/push_notification_integration_test.go#L17
	secretKeyTest = "188C879D394E09FDECC04606A126FAE2125FEABD24A2D12C6AC969AE1CEE2AEC"
)

func TestMain(m *testing.M) {
	cfg, err := New(instanceTest, secretKeyTest)
	if err != nil {
		log.Fatal(err)
	}

	c = cfg

	code := m.Run()

	os.Exit(code)
}

func TestNew(t *testing.T) {
	_, err := New("", "")

	if err == nil {
		log.Fatal("err must be not nil")
	}
}

func TestConfig_Send(t *testing.T) {
	type args struct {
		msg    Msg
		userID []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				msg:    map[string]any{"msg": "hello world"},
				userID: []string{"123123123123"},
			},
			wantErr: false,
		},

		{
			name: "empty users id",
			args: args{
				msg:    map[string]any{"msg": "asoyy"},
				userID: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Send(tt.args.msg, tt.args.userID...); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_GenerateToken(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		userID := "123123123"

		token, err := c.GenerateToken(userID)
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := token["token"]; !ok {
			t.Fatal("generated token is empty")
		}
	})

	t.Run("empty user id", func(t *testing.T) {
		_, err := c.GenerateToken("")
		if err == nil {
			t.Fatal("error must be not nil")
		}
	})
}
