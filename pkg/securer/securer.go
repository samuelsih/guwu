package securer

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"

	"golang.org/x/crypto/nacl/secretbox"
)

var secretKey [32]byte
var once sync.Once

var (
	ErrInvalidData = errors.New("invalid data")
	ErrInternal = errors.New("internal error")
)

func SetSecret(key [32]byte) {
	once.Do(func() {
		secretKey = key
	})
}

func Encrypt(input []byte) (string, error) {
	var nonce [24]byte

	_, err := rand.Read(nonce[:])
	if err != nil {
		return "", ErrInternal
	}

	box := secretbox.Seal(nonce[:], input, &nonce, &secretKey)

	return base64.RawURLEncoding.EncodeToString(box), nil
}

func Decrypt(input string) ([]byte, error) {
	box, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil || len(box) < 24 {
		return nil, ErrInvalidData 
	}

	var nonce [24]byte
	copy(nonce[:], box[:24])

	out, ok := secretbox.Open(nil, box[24:], &nonce, &secretKey)
	if ok {
		return out, nil
	}
	
	return nil, ErrInvalidData
}