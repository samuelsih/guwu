package securer

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"

	"github.com/samuelsih/guwu/pkg/errs"
	"golang.org/x/crypto/nacl/secretbox"
)

var secretKey [32]byte
var once sync.Once

var (
	ErrInvalidData = errors.New("invalid data")
	ErrInternal    = errors.New("internal error")
)

func SetSecret(key [32]byte) {
	once.Do(func() {
		secretKey = key
	})
}

func Encrypt(input []byte) (string, error) {
	const op = errs.Op("securer.Encrypt")
	var nonce [24]byte

	_, err := rand.Read(nonce[:])
	if err != nil {
		return "", errs.E(op, errs.KindUnexpected, err, "internal error")
	}

	box := secretbox.Seal(nonce[:], input, &nonce, &secretKey)

	return base64.RawURLEncoding.EncodeToString(box), nil
}

func Decrypt(input string) ([]byte, error) {
	const op = errs.Op("securer.Decrypt")

	box, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil || len(box) < 24 {
		return nil, errs.E(op, errs.KindBadRequest, err, "invalid data")
	}

	var nonce [24]byte
	copy(nonce[:], box[:24])

	out, ok := secretbox.Open(nil, box[24:], &nonce, &secretKey)
	if ok {
		return out, nil
	}

	return nil, errs.E(op, errs.KindBadRequest, err, "invalid data")
}
