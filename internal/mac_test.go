package internal

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encrypt_decrypt(t *testing.T) {
	plainText := []byte("samuelhotang02@gmail.com")

	cipherText, err := encrypt(plainText)

	assert.NoError(t, err)
	assert.NotEmpty(t, cipherText)

	text := base64.StdEncoding.EncodeToString(cipherText)

	textPlain, err := decrypt(text)

	assert.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Equal(t, textPlain, plainText)
}

func Test_GenerateMAC_ValidateMAC(t *testing.T) {
	result, err := GenerateMAC("samuelhotang02@gmail.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	
	validateText, err := ValidateMAC(result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	assert.Equal(t, "samuelhotang02@gmail.com", validateText.Email)
}
