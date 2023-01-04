package auth

import (
	"math/rand"
	"testing"
)

func Test_validEmail(t *testing.T) {
	t.Parallel()

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

	for _, tt := range tests {
		t.Run(tt.TestUsername, func(t *testing.T) {
			if tt.Result != validEmail(tt.Email) {
				t.Errorf("validEmail() = %v, want %v", validEmail(tt.Email), tt.Result)
			}
		})
	}
}

func Test_validUsername(t *testing.T) {
	t.Parallel()

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

	for _, tt := range tests {
		t.Run(tt.TestUsername, func(t *testing.T) {
			if tt.Result != validUsername(tt.Username) {
				t.Errorf("validEmail() = %v, want %v", validUsername(tt.Username), tt.Result)
			}
		})
	}
}

func Test_validPassword(t *testing.T) {
	t.Parallel()

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

	for _, tt := range tests {
		t.Run(tt.TestUsername, func(t *testing.T) {
			if tt.Result != validPassword(tt.Password) {
				t.Errorf("validPassword() = %v, want %v", validPassword(tt.Password), tt.Result)
			}
		})
	}
}

func Test_validateSignIn_Combine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    RegisterInput
		wantErr bool
	}{
		{"Empty", RegisterInput{}, true},
		{"Invalid Email", RegisterInput{Email: "invalid@example.com"}, true},
		{"Invalid Username", RegisterInput{Email: "samuelhotang02@gmail.com", Username: generateRandomStr(200)}, true},
		{"Invalid Password", RegisterInput{Email: "samuelhotang02@gmail.com", Username: generateRandomStr(20), Password: "aingmaung"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validAccount(tt.args.Username, tt.args.Email, tt.args.Password); (err != nil) != tt.wantErr {
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

