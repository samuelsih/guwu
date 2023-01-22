package business

import (
	"net/http"

	"github.com/samuelsih/guwu/pkg/errs"
	"github.com/samuelsih/guwu/pkg/logger"
)

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

func (res *CommonResponse) SetError(err error) {
	res.StatusCode = errs.GetKind(err)
	res.Msg = err.Error()

	logger.Err(err)
}

func (res *CommonResponse) RawError(statusCode int, msg string) {
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
