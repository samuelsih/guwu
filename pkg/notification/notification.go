package notification

import (
	"errors"
	pusher "github.com/pusher/push-notifications-go"
	"github.com/samuelsih/guwu/pkg/errs"
)

var (
	EmptyUserIDsErr = errors.New("empty users ID")
)

type Config struct {
	instance pusher.PushNotifications
}

type Msg map[string]any

func New(instanceID, secretKey string) (Config, error) {
	const op = errs.Op("notifications.New")

	instance, err := pusher.New(instanceID, secretKey)
	if err != nil {
		return Config{}, errs.E(op, errs.KindUnexpected, err, "unexpected internal server error")
	}

	return Config{
		instance: instance,
	}, nil
}

func (c Config) Send(msg Msg, userIDs ...string) error {
	const op = errs.Op("notifications.Send")

	if len(userIDs) == 0 || userIDs == nil {
		return errs.E(op, errs.KindUnexpected, EmptyUserIDsErr, "unexpected users")
	}

	req := map[string]any{
		"web": map[string]any{
			"notification": msg,
		},
	}

	_, err := c.instance.PublishToUsers(userIDs, req)
	if err != nil {
		return errs.E(op, errs.KindUnexpected, err, "cant send notification to user")
	}

	return nil
}

func (c Config) GenerateToken(userID string) (map[string]any, error) {
	const op = errs.Op("notifications.GenerateToken")

	token, err := c.instance.GenerateToken(userID)
	if err != nil {
		return nil, errs.E(op, errs.KindUnexpected, err, "can't authorize users")
	}

	return token, nil

}
