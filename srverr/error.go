package srverr

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Error struct {
	Code    string
	Message string
	Field   string
}

// Basic non-customer facing error message, e.g. log it and display generic error.
func NewError(msg string) *Error {
	return &Error{
		Message: msg,
	}
}

// Customer facing error message translatable using code.
func NewCodeError(msg, code string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// Customer facing error message for the provided field.
func NewCodeFieldError(msg, code, field string) *Error {
	return &Error{
		Code:    code,
		Field:   field,
		Message: msg,
	}
}

func (m *Error) Error() string {
	return m.Message
}

func (m *Error) FormatGQL(ctx context.Context) *gqlerror.Error {
	e := graphql.DefaultErrorPresenter(ctx, m)

	var se *Error
	if errors.As(e, &se) {
		e.Message = se.Message
		e.Extensions = make(map[string]any)

		if len(se.Code) > 0 {
			e.Extensions["code"] = se.Code
		}
		if len(se.Field) > 0 {
			e.Extensions["field"] = se.Field
		}
	}

	return e
}
