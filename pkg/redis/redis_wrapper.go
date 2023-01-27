package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/rueian/rueidis"
	"github.com/samuelsih/guwu/pkg/errs"
)

var (
	ErrInternal         = errors.New("internal error")
	ErrInvalidUnmarshal = errors.New("can't unmarshal")
	ErrInvalidMarshal   = errors.New("can't marshal")
	ErrUnknownKey       = errors.New("unknown input")
)

type Client struct {
	Pool   rueidis.Client
}

func NewClient(db rueidis.Client) *Client {
	return &Client{
		Pool: db,
	}
}

func (r *Client) Get(ctx context.Context, key string) (string, error) {
	const op = errs.Op("redis_wrapper.Get")
	result, err := r.Pool.Do(ctx, r.Pool.B().Get().Key(key).Build()).ToString()

	if err != nil && !(err.Error() == `redis nil message` || err.Error() == `redis: nil`) {
		return "", errs.E(op, errs.KindUnexpected, err, "internal error")
	}

	return strings.Trim(result, `"`), nil
}

func (r *Client) Set(ctx context.Context, key, value string, time int64) error {
	const op = errs.Op("redis_wrapper.Set")

	err := r.Pool.Do(ctx, r.Pool.B().Setex().Key(key).Seconds(time).Value(value).Build()).Error()

	if err != nil {
		return errs.E(op, errs.KindUnexpected, err, "internal error")
	}

	return nil
}

func (r *Client) GetJSON(ctx context.Context, key string, dst any) error {
	const op = errs.Op("redis_wrapper.GetJSON")

	result, err := r.Pool.Do(ctx, r.Pool.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if !(err.Error() == `redis nil message` || err.Error() == `redis: nil`) {
			return errs.E(op, errs.KindUnexpected, err, "internal error")
		}

		return errs.E(op, errs.KindBadRequest, err, "unknown input")
	}

	err = json.Unmarshal([]byte(result), dst)
	if err != nil {
		return errs.E(op, errs.KindBadRequest, err, "cannot unmarshal")
	}

	return nil
}

func (r *Client) SetJSON(ctx context.Context, key string, value any, time int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrInvalidUnmarshal
	}

	err = r.Pool.Do(ctx, r.Pool.B().Setex().Key(key).Seconds(time).Value(string(data)).Build()).Error()
	if err != nil {
		return ErrInternal
	}

	return nil
}

func (r *Client) Destroy(ctx context.Context, key string) error {
	deleted, err := r.Pool.Do(ctx, r.Pool.B().Del().Key(key).Build()).ToInt64()
	if err != nil {
		return ErrInternal
	}

	if deleted == 0 {
		return ErrUnknownKey
	}

	return nil
}
