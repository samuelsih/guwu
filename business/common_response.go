package business

import "net/http"

var _ CommonOutput = (*CommonResponse)(nil)

type CommonResponse struct {
	StatusCode int    `json:"code,omitempty"`
	Msg        string `json:"message,omitempty"`
	SessionID  string `json:"-"`
}

func (res *CommonResponse) SetError(statusCode int, msg string) {
	res.StatusCode = statusCode
	res.Msg = msg
}

func (res *CommonResponse) SetOK() {
	res.StatusCode = http.StatusOK
	res.Msg = "OK"
}

func (res CommonResponse) CommonRes() *CommonResponse {
	return &res
}

type CommonOutput interface {
	CommonRes() *CommonResponse
}
