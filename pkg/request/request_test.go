package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "embed"
)

//go:embed sample.txt
var sampleText string

func Test_DecodeSuccess(t *testing.T) {
	type someData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	data := someData{Name: "bob", Age: 14}

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	body := strings.NewReader(string(b))

	r := httptest.NewRequest(http.MethodGet, "/", body)
	w := httptest.NewRecorder()

	var dst someData

	err = Decode(w, r, &dst)
	if err != nil {
		t.Fatalf("Unable to decode: %v", err)
	}
}

func Test_EmptyBody(t *testing.T) {
	body := strings.NewReader("")
	var data map[string]any

	r := httptest.NewRequest(http.MethodGet, "/", body)
	w := httptest.NewRecorder()

	err := Decode(w, r, &data)
	if err.Error() != "body must not be empty" {
		t.Fatalf("expected %v, got %v", "body must not be empty", err)
	}
}

func Test_SyntaxError(t *testing.T) {
	body := strings.NewReader(
		`{
			123: 123
		 }`,
	)

	var data map[string]any

	r := httptest.NewRequest(http.MethodGet, "/", body)
	w := httptest.NewRecorder()

	var syntaxError *json.SyntaxError

	err := Decode(w, r, &data)
	if errors.As(err, &syntaxError) {
		t.Fatalf("expected %v, got %v", syntaxError.Error(), err)
	}
}

func Test_UnmarshalTypeError(t *testing.T) {
	body := strings.NewReader(
		`{
			"name": false
		 }`,
	)

	type A struct {
		Name string `json:"name"`
	}

	var data A

	r := httptest.NewRequest(http.MethodGet, "/", body)
	w := httptest.NewRecorder()

	err := Decode(w, r, &data)

	if !strings.HasPrefix(err.Error(), "body contains incorrect JSON") {
		t.Fatalf("expected %v, got %v", "body contains incorrect JSON", err)
	}
}

func Test_UnknownFields(t *testing.T) {
	body := strings.NewReader(`
		{
			"names": "Bob"
		}
	`)

	type A struct {
		Name string `json:"name"`
	}

	var data A

	r := httptest.NewRequest(http.MethodGet, "/", body)
	w := httptest.NewRecorder()

	err := Decode(w, r, &data)
	if !strings.HasPrefix(err.Error(), "body contains unknown key ") {
		t.Fatalf("expected %v, got %v", "body contains unknown key", err)
	}
}

func Test_InvalidUnmarshalError(t *testing.T) {
	body := strings.NewReader(`
	{
		"name": "foo"
	}`)

	var data map[string]any

	r := httptest.NewRequest(http.MethodGet, "/foo", body)
	w := httptest.NewRecorder()

	err := Decode(w, r, data)

	if err.Error() != "invalid request body" {
		t.Fatalf("expected %v, got %v", "invalid request body", err)
	}
}

func Test_BodyTooLarge(t *testing.T) {
	in := fmt.Sprintf(`{"foo": "%s"}`, sampleText)

	body := strings.NewReader(in)

	var data map[string]any

	r := httptest.NewRequest(http.MethodGet, "/foo", body)
	w := httptest.NewRecorder()

	err := Decode(w, r, &data)

	if !strings.HasPrefix(err.Error(), "body must not be larger than") {
		t.Fatalf("expected %v, got %v", "request body too large", err)
	}
}

func Test_SingleJSONValue(t *testing.T) {
	in := `
		{
			"name": "Foo"
		}
		{
			"age": "bar"
		}
	`

	body := strings.NewReader(in)

	r := httptest.NewRequest(http.MethodGet, "/foo", body)
	w := httptest.NewRecorder()

	var data map[string]any

	err := Decode(w, r, &data)
	t.Log(err)

	if err.Error() != "body must only contain a single JSON value" {
		t.Fatalf("expected %v, got %v", "body must only contain a single JSON value", err)
	}
}
