package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/rueian/rueidis"
)

var (
	ErrInternal = errors.New("internal error")
	ErrInvalidUnmarshal = errors.New("can't unmarshal")
	ErrInvalidMarshal = errors.New("can't marshal")
	ErrUnknownKey = errors.New("unknown input'")
)

type Client struct {
	Pool   rueidis.Client
	Prefix string
}

func NewClient(db rueidis.Client, prefix string) *Client {
	return &Client{
		Pool: db,
	}
}

func (r *Client) Get(ctx context.Context, key string) (string, error) {
	result, err := r.Pool.Do(ctx, r.Pool.B().Get().Key(r.Prefix + key).Build()).ToString()

	if err != nil && !(err.Error() == `redis nil message` || err.Error() == `redis: nil`) {
		return "", ErrInternal
	}

	return strings.Trim(result, `"`), nil
}

func (r *Client) Set(ctx context.Context, key, value string, time int64) error {
	err := r.Pool.Do(ctx, r.Pool.B().Setex().Key(r.Prefix+key).Seconds(time).Value(value).Build()).Error()

	if err != nil {
		return ErrInternal
	}

	return nil
}

func (r *Client) GetJSON(ctx context.Context, key string, dst any) error {
	result, err := r.Pool.Do(ctx, r.Pool.B().Get().Key(r.Prefix + key).Build()).ToString()

	if err != nil {
		return ErrInternal
	}

	err = json.Unmarshal([]byte(result), dst)
	if err != nil {
		return ErrInvalidUnmarshal
	}

	return nil
}

func (r *Client) SetJSON(ctx context.Context, key string, value any, time int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrInvalidUnmarshal
	}

	err = r.Pool.Do(ctx, r.Pool.B().Setex().Key(r.Prefix+key).Seconds(time).Value(string(data)).Build()).Error()
	if err != nil {
		return ErrInternal
	}

	return nil
}

func (r *Client) Destroy(ctx context.Context, key string) error {
	deleted, err := r.Pool.Do(ctx, r.Pool.B().Del().Key(r.Prefix+key).Build()).ToBool()
	if err != nil {
		return ErrInternal
	}

	if !deleted {
		return ErrUnknownKey
	}

	return nil
}
