package business

import "net/http"

// type check interface
var _ CommonOutput = (*CommonResponse)(nil)

type CommonInputMatcher interface{}
type CommonInput struct {
	SessionID  string
	URLParam   []string
	QueryParam []string
}

type CommonResponse struct {
	StatusCode    int    `json:"code,omitempty"`
	Msg           string `json:"message,omitempty"`
	SessionID     string `json:"-"`
	SessionMaxAge int    `json:"-"`
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
