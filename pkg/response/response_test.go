package response

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_JSON(t *testing.T) {
	data := struct {
		name  string
		value string
	}{
		name:  "foo",
		value: "bar",
	}

	w := httptest.NewRecorder()

	err := JSON(w, 200, data)
	if err != nil {
		t.Fatalf("expect %v, got %v", nil, err)
	}
}

func Test_JSONWithHeaders(t *testing.T) {
	data := struct {
		name string
		age  int
	}{
		name: "test",
		age:  14,
	}

	var header http.Header = make(http.Header)
	header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	err := JSONWithHeaders(w, 200, data, header)
	if err != nil {
		t.Fatalf("expect %v, got %v", nil, err)
	}
}

func Test_JSONWithHeaders_MustFail(t *testing.T) {
	var typeErr *json.UnsupportedTypeError
	var valErr *json.UnsupportedValueError

	data := func() {}

	w := httptest.NewRecorder()
	var header http.Header = make(http.Header)
	header.Set("Content-Type", "application/json")

	err := JSONWithHeaders(w, 200, data, header)

	if !errors.As(err, &typeErr) {
		t.Fatalf("expect %v, got %v", nil, err)
	}

	data2 := math.Inf(1)
	err = JSONWithHeaders(w, 200, data2, header)
	if !errors.As(err, &valErr) {
		t.Fatalf("expect %v, got %v", nil, err)
	}
}
