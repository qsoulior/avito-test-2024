package service

type ErrorType string

var (
	ErrorTypeInvalid      ErrorType = "invalid"      // 400
	ErrorTypeUnauthorized ErrorType = "unauthorized" // 401
	ErrorTypeForbidden    ErrorType = "forbidden"    // 403
	ErrorTypeNotExist     ErrorType = "not exist"    // 404
	ErrorTypeInternal     ErrorType = "internal"     // 500
)

type Error struct {
	msg   string
	etype ErrorType
	err   error
}

func NewTypedError(msg string, etype ErrorType, err error) error { return &Error{msg, etype, err} }

func (e *Error) Error() string { return e.msg }

func (e *Error) Unwrap() error { return e.err }

func (e *Error) Type() ErrorType { return e.etype }
