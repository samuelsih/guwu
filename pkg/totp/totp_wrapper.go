package totp

import (
	"time"

	"github.com/xlzd/gotp"
)

type TOTP struct {
	t *gotp.TOTP
} 

func New(key string, duration, digits int) TOTP {
	return TOTP{
		t: gotp.NewTOTP(key, digits, duration, nil),
	}
}

func (totp TOTP) Generate() (string, int64) {
	return totp.t.NowWithExpiration()
}

func(totp TOTP) Verify(otpCode string) bool {
	return totp.t.Verify(otpCode, time.Now().Unix())
}