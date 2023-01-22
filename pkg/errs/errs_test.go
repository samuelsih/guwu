package errs

import (
	"errors"
	"testing"
)

func TestErrs(t *testing.T) {
	err := fnA()
	err = E("main", GetKind(err), err, "unexpected")

	output := Ops(err.(*Error))
	expected := []Op{"main", "fnA", "fnB"}

	if !equal(output, expected) {
		t.Fatalf("expected %v, got %v", expected, output)
	}
}

func fnA() error {
	err := fnB()
	return E("fnA", GetKind(err), err, "client error from fnB")
}

func fnB() error {
	e := errors.New("error from fnB")
	return E("fnB", 500, e, "client error from fnB")
}

func equal(a, b []Op) bool {
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
