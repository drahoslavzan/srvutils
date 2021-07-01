package gqlerrs

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func UnauthorizedError() *gqlerror.Error {
	return &gqlerror.Error{
		Extensions: map[string]interface{}{
			"ecode": "auth/unauthorized",
		},
	}
}
