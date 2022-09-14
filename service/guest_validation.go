package service

import (
	"errors"
	"net"
	"net/mail"
	"strings"
	"unicode"

	"github.com/rs/zerolog/log"
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

	errNameRequired  = errors.New("name is required")
	errNameMaxLength = errors.New("name length must be lower than 100 characters")
)

func validateSignIn(data GuestRegisterIn) error {
	if err := validateEmail(data.Email); err != nil {
		log.Debug().Stack().Err(err).Str("place", "validate email")
		return err
	}

	if err := validateName(data.Name); err != nil {
		log.Debug().Stack().Err(err).Str("place", "validate name")
		return err
	}

	if err := validatePassword(data.Password); err != nil {
		log.Debug().Stack().Err(err).Str("place", "validate password")
		return err
	}

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return errEmailRequired
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errInvalidEmail
	}

	if !strings.Contains(email, "@") {
		return errInvalidEmail
	}

	domparts := strings.Split(email, "@")

	_, err = net.LookupMX(domparts[1])
	if err != nil {
		return errDomainNotFound
	}

	return nil
}

func validateName(name string) error {
	if name == "" {
		return errNameRequired
	}

	if len(name) > 99 {
		return errNameMaxLength
	}

	return nil
}

func validatePassword(password string) error {
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