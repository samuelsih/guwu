package service

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateEmail(t *testing.T) {
	tests := []struct {
		TestName string
		Email string
		Result error
	}{
		{"Empty", "", errEmailRequired},
		{"Invalid", "John Doe", errInvalidEmail},
		{"Not_Contain_@", "bujanginam123", errInvalidEmail},
		{"Valid_But_Non_Existent", "asem@ididntexists.123", errDomainNotFound},
		{"Valid_And_Exists", "samuelhotang02@gmail.com", nil},
	}

	for _, item := range tests {
		t.Run(item.TestName, func(t *testing.T) {
			assert.Equal(t, item.Result, validateEmail(item.Email))
		})
	}
}

func Test_validateName(t *testing.T) {
	tests := []struct {
		TestName string
		Name string
		Result error
	}{
		{"Empty", "", errNameRequired},
		{"More_Than_100_Chars", generateRandomStr(101), errNameMaxLength},
		{"Valid", "Dadang Maddog", nil},
		{"Valid_With_Random_String", generateRandomStr(50), nil},
	}

	for _, item := range tests {
		t.Run(item.TestName, func(t *testing.T) {
			assert.Equal(t, item.Result, validateName(item.Name))
		})
	}
}

func Test_validatePassword(t *testing.T) {
	tests := []struct {
		TestName string
		Password string
		Result error
	}{
		{"Empty", "", errPasswordRequired},
		{"Lower_Than_8", generateRandomStr(7), errPasswordLength},
		{"Not_Contain_Lowercase", "HAYANGULIN123", errPasswordLowerChar},
		{"Not_Contain_Uppercase", "hayangulin123", errPasswordUpperChar},
		{"Not_Contain_Symbol", "Hayangulin123", errPasswordSymbolChar},
		{"Not_Contain_Number", "Hayangulin!!", errPasswordNumChar},
		{"Valid", "Hayangulin123!!", nil},
	}

	for _, item := range tests {
		t.Run(item.TestName, func(t *testing.T) {
			assert.Equal(t, item.Result, validatePassword(item.Password))
		})
	}
}

func generateRandomStr(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}