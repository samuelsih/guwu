package securer

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	generateKey()

	plainText := "some_unique_id_for_session_or_apapun"

	encrypted, err := Encrypt([]byte(plainText))
	if err != nil {
		t.Fatalf("Err should nil, got %v", err)
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Err should nil, got %v", err)
	}

	if !bytes.Equal(decrypted, []byte(plainText)) {
		t.Fatalf("result are not same")
	}
}

func TestEncodeDecodeFails(t *testing.T) {
	generateKey()

	plainText := "some_unique_id_for_session_or_apapun"

	encrypted, err := Encrypt([]byte(plainText))
	if err != nil {
		t.Fatalf("Encrypted: err should not nil, got %v", err)
	}

	// try to change the encrypted data
	r := []rune(encrypted)
	r[1] = rune(r[10])

	_, err = Decrypt(string(r))
	if err == nil {
		t.Fatalf("Decrpyt: err should not nil, got nil: %v", err)
	}
}

func generateKey() {
	b, _ := hex.DecodeString("9732070617373776f726420746f206120736563726574")

	copy(secretKey[:], b)
}