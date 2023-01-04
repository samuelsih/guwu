package auth

import (
	"errors"
	"net"
	"net/mail"
	"strings"
	"unicode"
)

var (
	errEmailRequired  = errors.New("email is required")
	errInvalidEmail   = errors.New("email is invalid")
	errDomainNotFound = errors.New("email domain not found")

	errPasswordRequired   = errors.New("password is required")
	errPasswordLength     = errors.New("password must be more than 8 characters")
	errPasswordLowerChar  = errors.New("password must contains lower char")
	errPasswordUpperChar  = errors.New("password must contains upper char")
	errPasswordSymbolChar = errors.New("password must contains symbol char")
	errPasswordNumChar    = errors.New("password must contains number char")

	errUsernameRequired  = errors.New("username is required")
	errUsernameMaxLength = errors.New("username length must be lower than 50 characters")
)

func validAccount(username, email, password string) error {
	if err := validEmail(email); err != nil {
		return err
	}

	if err := validUsername(username); err != nil {
		return err
	}

	if err := validPassword(password); err != nil {
		return err
	}

	return nil
}

func validEmail(email string) error {
	if email == "" {
		return errEmailRequired
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errInvalidEmail
	}

	domparts := strings.Split(email, "@")

	_, err = net.LookupMX(domparts[1])
	if err != nil {
		return errDomainNotFound
	}

	return nil
}

func validUsername(name string) error {
	if name == "" {
		return errUsernameRequired
	}

	if len(name) > 99 {
		return errUsernameMaxLength
	}

	return nil
}

func validPassword(password string) error {
	if password == "" {
		return errPasswordRequired
	}

	isMoreThan8 := len(password) > 8
	if !isMoreThan8 {
		return errPasswordLength
	}

	var isLower, isUpper, isSymbol, isNumber bool

	for _, p := range password {
		if !isLower && unicode.IsLower(p) {
			isLower = true
		}

		if !isUpper && unicode.IsUpper(p) {
			isUpper = true
		}

		if !isSymbol && (unicode.IsSymbol(p) || unicode.IsPunct(p)) {
			isSymbol = true
		}

		if !isNumber && unicode.IsNumber(p) {
			isNumber = true
		}
	}

	if !isLower {
		return errPasswordLowerChar
	}

	if !isUpper {
		return errPasswordUpperChar
	}

	if !isSymbol {
		return errPasswordSymbolChar
	}

	if !isNumber {
		return errPasswordNumChar
	}

	return nil
}
