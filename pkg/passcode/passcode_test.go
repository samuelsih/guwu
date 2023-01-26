package passcode

import "testing"

func TestGenerateUnique(t *testing.T) {
	var code string
	
	for i := 0; i < 10; i++ {
		generatedPasscode := Generate(6)
		if code == generatedPasscode {
			t.Fatal("code not unique")
		}
	}
}