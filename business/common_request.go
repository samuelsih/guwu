package business

var _ CommonInput = (*CommonRequest)(nil)

type CommonRequest struct {
	SessionID string `json:"-"`
}

func (req CommonRequest) CommonReq() *CommonRequest {
	return &req
}

type CommonInput interface {
	CommonReq() *CommonRequest
}
