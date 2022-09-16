package service

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateEmail(t *testing.T) {
	tests := []struct {
		TestUsername string
		Email        string
		Result       error
	}{
		{"Empty", "", errEmailRequired},
		{"Invalid_Contain_Space", "email at gmail.com", errInvalidEmail},
		{"Not_Contain_@", "bujanginam123", errInvalidEmail},
		{"Valid_But_Non_Existent", "asem@ididntexists.123", errDomainNotFound},
		{"Valid_And_Exists", "samuelhotang02@gmail.com", nil},
	}

	for _, item := range tests {
		t.Run(item.TestUsername, func(t *testing.T) {
			assert.Equal(t, item.Result, validateEmail(item.Email))
		})
	}
}

func Test_validateUsername(t *testing.T) {
	tests := []struct {
		TestUsername string
		Username     string
		Result       error
	}{
		{"Empty", "", errUsernameRequired},
		{"More_Than_100_Chars", generateRandomStr(101), errUsernameMaxLength},
		{"Valid", "Dadang Maddog", nil},
		{"Valid_With_Random_String", generateRandomStr(50), nil},
	}

	for _, item := range tests {
		t.Run(item.TestUsername, func(t *testing.T) {
			assert.Equal(t, item.Result, validateUsername(item.Username))
		})
	}
}

func Test_validatePassword(t *testing.T) {
	tests := []struct {
		TestUsername string
		Password     string
		Result       error
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
		t.Run(item.TestUsername, func(t *testing.T) {
			assert.Equal(t, item.Result, validatePassword(item.Password))
		})
	}
}

func Test_validateSignIn_Combine(t *testing.T) {
	tests := []struct {
		name    string
		args    GuestRegisterIn
		wantErr bool
	}{
		{"Empty", GuestRegisterIn{}, true},
		{"Invalid Email", GuestRegisterIn{Email: "invalid@example.com"}, true},
		{"Invalid Username", GuestRegisterIn{Email: "samuelhotang02@gmail.com", Username: generateRandomStr(200)}, true},
		{"Invalid Password", GuestRegisterIn{Email: "samuelhotang02@gmail.com", Username: generateRandomStr(20), Password: "aingmaung"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSignIn(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("validateSignIn() error = %v, wantErr %v", err, tt.wantErr)
			}
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


