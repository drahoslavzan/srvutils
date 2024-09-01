package gqlerrs

type Error struct {
	Code    string
	Message string
	Field   string
}

func NewError(msg string) *Error {
	return &Error{
		Message: msg,
	}
}

func NewCodeError(msg, code string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

func NewFieldError(msg, field string) *Error {
	return &Error{
		Message: msg,
		Field:   field,
	}
}

func (m *Error) Error() string {
	return m.Message
}
