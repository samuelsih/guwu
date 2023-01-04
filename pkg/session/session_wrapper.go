package session

import "errors"

var (
	UnknownSessionID = errors.New("unknown session id")
	InternalErr = errors.New("can't create session for now, try again later")
)