package model

import "net/http"

// Success
var (
	OK = http.StatusOK
	Created = http.StatusCreated
	NoContent = http.StatusNoContent
)


// Client Error
var (
	BadRequest = http.StatusBadRequest
	Unauthorized = http.StatusUnauthorized
	NotFound = http.StatusNotFound
)

// Server Error
var (
	InternalServerError = http.StatusInternalServerError
)