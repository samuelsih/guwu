package securer

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"log"
	"math/big"
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
	toDecrypt := randString(encrypted)

	_, err = Decrypt(toDecrypt)
	if err == nil {
		t.Fatalf("Decrpyt: err should not nil, got nil: %v", err)
	}
}

func generateKey() {
	b, _ := hex.DecodeString("9732070617373776f726420746f206120736563726574")

	copy(secretKey[:], b)
}

func randString(str string) string {
	r := []rune(str)
	l := len(r) - 1

	n, err := rand.Int(rand.Reader, big.NewInt(int64(l)))
	if err != nil {
		log.Fatalf("randString: %v", err)
	}

	r = append(r, rune(int(n.Int64())))

	return string(r)
}