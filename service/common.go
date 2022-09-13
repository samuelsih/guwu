package service

import "net/http"

type CommonResponse struct {
	StatusCode int    `json:"code,omitempty"`
	Msg        string `json:"message,omitempty"`
}

func (o *CommonResponse) SetError(statusCode int, errMsg string) {
	o.StatusCode = statusCode
	o.Msg = errMsg
}

func (o *CommonResponse) SetOK() {
	o.StatusCode = http.StatusOK
	o.Msg = "OK"
}

func (o *CommonResponse) SetCreated() {
	o.StatusCode = http.StatusCreated
	o.Msg = "Created"
}

