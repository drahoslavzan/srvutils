package server

import (
	"context"

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
func (m *Error) WithCode(code string) *Error {
	m.Code = code
	return m
}

// Customer facing error message for the provided field.
func (m *Error) OnField(field string) *Error {
	m.Field = field
	return m
}

func (m *Error) Error() string {
	return m.Message
}

func (m *Error) FormatGQL(ctx context.Context) *gqlerror.Error {
	e := graphql.DefaultErrorPresenter(ctx, m)

	e.Message = m.Message
	e.Extensions = make(map[string]any)
	if len(m.Code) > 0 {
		e.Extensions["code"] = m.Code
	}
	if len(m.Field) > 0 {
		e.Extensions["field"] = m.Field
	}

	return e
}
