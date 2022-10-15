package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	passphrase = os.Getenv("MAC_KEY")
	iv = make([]byte, 16)
)

type Payload struct {
	Email string `json:"email"`
	Duration time.Time `json:"duration"`
}

func GenerateMAC(email string) (string, error) {
	plainText, err := json.Marshal(Payload{Email: email, Duration: time.Now().Local().Add(time.Minute * 10)})
	if err != nil {
		return "", err
	}

	cipherText, err := encrypt(plainText)
	if err != nil {
		return "", err
	}

	cipherString := base64.StdEncoding.EncodeToString(cipherText)
	return cipherString, nil
}


func ValidateMAC(mac string) (Payload, error) {
	plainText, err := decrypt(mac)
	if err != nil {
		return Payload{}, err
	}

	var payload Payload
	if err := json.Unmarshal(plainText, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}

func encrypt(plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))

	mode.XORKeyStream(cipherText, plainText)

	return cipherText, nil
}


func decrypt(text string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	} 
	
	mode := cipher.NewCFBDecrypter(block, iv)
	
	plainText := make([]byte, len(cipherText)) 
	
	mode.XORKeyStream(plainText, cipherText)

	return plainText, nil
}

