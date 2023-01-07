package errs

var _ error = (*Error)(nil)

type Error struct {
	Op   string
	Kind int
	Err  error
	Raw  any
}

func (e Error) Error() string {
	return e.Err.Error()
}

// Generics dont support error type, so any is used
func E(args ...any) error {
	e := &Error{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.Op = arg

		case error:
			e.Err = arg

		default:
			e.Raw = arg
		}
	}

	return e
}

func Ops(e *Error) []string {
	res := []string{e.Op}

	for {
		subErr, ok := e.Err.(*Error)
		if !ok {
			return res
		}

		res = append(res, subErr.Op)

		e = subErr
	}
}