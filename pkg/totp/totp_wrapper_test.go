package totp

import (
	"testing"
	"time"
)

func TestTOTP(t *testing.T) {
	totp := New("4S62BZNFXXSZLCRO", (5 * 60), 6)

	code, _ := totp.Generate()

	if !totp.Verify(code) {
		t.Fatal("duration error")
	} 
}

func TestTOTP_Fail(t *testing.T) {
	totp := New("4S62BZNFXXSZLCRO", (1), 6)
	code, _ := totp.Generate()

	time.Sleep(2 * time.Second)

	if totp.Verify(code) {
		t.Fatal("duration fail error")
	} 
}