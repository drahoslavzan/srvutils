package gqlerrs

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

type srvError struct {
	code   string
	fields map[string]interface{}
}

type SrvError interface {
	Code() string
	Fields() map[string]interface{}
	Graph() *gqlerror.Error
}

func MakeSrvError(code string, fields map[string]interface{}) SrvError {
	return &srvError{
		code:   code,
		fields: fields,
	}
}

func (m *srvError) Code() string {
	return m.code
}

func (m *srvError) Fields() map[string]interface{} {
	return m.fields
}

func (m *srvError) Error() string {
	if len(m.code) > 0 {
		if len(m.fields) > 0 {
			return fmt.Sprintf("%s: %v", m.code, m.fields)
		}
		return fmt.Sprintf("%s", m.code)
	}
	return fmt.Sprintf("%v", m.fields)
}

func (m *srvError) Graph() *gqlerror.Error {
	return &gqlerror.Error{
		Extensions: map[string]interface{}{
			"ecode":  m.Code(),
			"fields": m.Fields(),
		},
	}
}
