package service

import (
	"net/http"

	"github.com/samuelsih/guwu/model"
)

type CommonRequest struct {
	Token string `json:"-"`
	UserSession model.Session `json:"-"`
}
type CommonResponse struct {
	StatusCode int    `json:"code,omitempty"`
	Msg        string `json:"message,omitempty"`
	SessionID string `json:""`
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

type CommonOutput interface {
	CommonRes() *CommonResponse
}

func (o CommonResponse) CommonRes() *CommonResponse {
	return &o
}

type CommonInput interface {
	CommonReq() *CommonRequest
}

func (o CommonRequest) CommonReq() *CommonRequest {
	return &o
}
