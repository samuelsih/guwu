package totp

import "github.com/xlzd/gotp"

type TOTP struct {
	t *gotp.TOTP
} 

func New(key string, duration, digits int) TOTP {
	return TOTP{
		t: gotp.NewTOTP(key, digits, duration, nil),
	}
}

func (tt TOTP) Generate() (string, int64) {
	return tt.t.NowWithExpiration()
}

func(tt TOTP) Verify(otpCode string, time int64) bool {
	return tt.t.Verify(otpCode, time)
}