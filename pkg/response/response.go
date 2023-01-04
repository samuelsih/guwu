package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data any) error {
	return JSONWithHeaders(w, status, data, nil)
}

func JSONWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	for key, value := range headers {
		w.Header()[key] = value
	}

	var valueErr *json.UnsupportedValueError
	var typeErr *json.UnsupportedTypeError

	enc := json.NewEncoder(w)

	err := enc.Encode(data)
	if err != nil {
		switch {
		case errors.As(err, &valueErr):
			return fmt.Errorf("unsupported value for this response %w", err)

		case errors.As(err, &typeErr):
			return fmt.Errorf("unsupported type for this response %w", err)

		default:
			return fmt.Errorf("can't encode the response: %w", err)
		}
	}

	return nil
}
