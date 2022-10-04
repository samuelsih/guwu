package model

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/go-redis/redis/v8"
)

var (
	day = 24 * 60 * 3600
	mx sync.RWMutex
)

type Session struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type SessionDeps struct {
	Conn *redis.Client
}

func (u *SessionDeps) Save(ctx context.Context, data Session) (string, error) {
	mx.RLock()
	defer mx.RUnlock()
	
	buf, err := sonic.Marshal(&data)
	if err != nil {
		return "", err
	}

	value := string(buf)

	sessionID, err := generateSessID()
	if err != nil {
		return "", err
	}
	
	cmd := u.Conn.Set(ctx, sessionID, value, time.Duration(day) * time.Second)

	return sessionID, cmd.Err()
}

func (u *SessionDeps) Get(ctx context.Context, key string) (Session, error) {
	var userSession Session

	userFromDB, err := u.Conn.Get(ctx, key).Result()
	if err != nil {
		return userSession, err
	}

	r := bytes.NewReader([]byte(userFromDB))

	err = decoder.NewStreamDecoder(r).Decode(&userSession)
	if err != nil {
		return userSession, err
	}

	return userSession, err
}

func (u *SessionDeps) Delete(ctx context.Context, key string) error {
	if len(key) == 0 {
		return errors.New(`user not found`)
	}

	cmd := u.Conn.Del(ctx, key)

	if cmd.Err() != nil {
		return errors.New(`user not found`)
	}

	return nil
}

func generateSessID() (string, error) {
	b := make([]byte, 15)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil

}
