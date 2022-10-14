package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"time"
)

var (
	secretKey = []byte(os.Getenv("MAC_KEY"))
	iv =  make([]byte, 16)
)

const blockSize = 32

type Payload struct {
	Email string `json:"email"`
	Duration time.Time `json:"duration"`
}


func GenerateMAC(email string) (string, error) {
	phrase, err := json.Marshal(Payload{Email: email, Duration: time.Now().Local().Add(time.Minute * 10)})
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(secretKey)
	key := h.Sum(nil)

	cipherText, err := encrypt(key, phrase)
	if err != nil {
		return "", err
	}

	cipherString := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(cipherText)
	return cipherString, nil
}


func ValidateMAC(mac string) (Payload, error) {
	h := sha256.New()
	h.Write([]byte(secretKey))
	key := h.Sum(nil)

	plainText, err := decrypt(key, []byte(mac))
	if err != nil {
		return Payload{}, err
	}

	var payload Payload
	if err := json.Unmarshal(plainText, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}

func encrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	plainText = pkcs7pad256(plainText)
	cipherText := make([]byte, len(plainText))

	mode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}


func decrypt(key, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))

	mode.CryptBlocks(plainText, cipherText)

	return pkcs7strip256(plainText)
}

// add paddding if data length is < 32 bytes
func pkcs7pad256(data []byte) []byte {
	dataLen := len(data)
	paddingLen := blockSize % dataLen
	padding := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)
	return append(data, padding...)
}

// check padding 
func pkcs7strip256(data []byte) ([]byte, error) {
	dataLen := len(data)

	if dataLen == 0 {
		return nil, errors.New("data is empty")
	}

	if dataLen % blockSize != 0 {
		return nil, errors.New("data is not block-aligned with blocksize")
	}

	paddingLen := int(data[dataLen - 1])
	ref := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)

	if paddingLen > blockSize || paddingLen == 0 || !bytes.HasSuffix(data, ref) {
		return nil, errors.New("invalid padding")
	}

	return data[:dataLen - paddingLen], nil
}