package errs

import (
	"errors"
	"testing"
)

func TestErrs(t *testing.T) {
	err := fnA()
	err = E("main", errors.New("dari main"), err) 	

	output := Ops(err.(*Error))
	expected := []string{"main", "fnA", "fnB"}

	if !equal(output, expected) {
		t.Fatalf("expected %v, got %v", expected, output)
	}
}

func fnA() error {
	err := fnB()

	return E("fnA", errors.New("error dari fnA"), err)
}

func fnB() error {
	return E("fnB", errors.New("error dari fnB"))
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}