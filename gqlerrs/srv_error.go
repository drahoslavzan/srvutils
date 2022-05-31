package gqlerrs

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

type srvError struct {
	code   string
	fields map[string]any
}

type SrvError interface {
	Code() string
	Fields() map[string]any
	Plain() error
	Graph() *gqlerror.Error
}

func MakeSrvError(code string, fields map[string]any) SrvError {
	return &srvError{
		code:   code,
		fields: fields,
	}
}

func (m *srvError) Code() string {
	return m.code
}

func (m *srvError) Fields() map[string]any {
	return m.fields
}

func (m *srvError) Plain() error {
	if len(m.code) > 0 {
		if len(m.fields) > 0 {
			return fmt.Errorf("%s: %v", m.code, m.fields)
		}
		return fmt.Errorf("%s", m.code)
	}
	return fmt.Errorf("%v", m.fields)
}

func (m *srvError) Graph() *gqlerror.Error {
	return &gqlerror.Error{
		Extensions: map[string]any{
			"ecode":  m.Code(),
			"fields": m.Fields(),
		},
	}
}
