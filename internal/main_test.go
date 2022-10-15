package internal

import "testing"

func TestMain(m *testing.M) {
	if passphrase == "" {
		passphrase = "SODNOLSOWXGGYJNZKEMCYBQHUWIAMWTI"
	}
}