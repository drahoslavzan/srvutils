package validation

import (
	"net"
	"regexp"
	"strings"
)

var basicEmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsEmailValid checks if the given email address is valid by verifying:
// 1. The email format using the provided regex pattern (or the default if nil).
// 2. The email length is between 3 and 254 characters.
// 3. The domain has valid MX (Mail Exchange) records.
// It returns true if the email is valid, false otherwise.
func IsEmailValid(email string, regex *regexp.Regexp) bool {
	if regex == nil {
		regex = basicEmailRegex
	}

	if len(email) < 3 || len(email) > 254 {
		return false
	}

	if !regex.MatchString(email) {
		return false
	}

	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return false
	}

	domain := email[atIndex+1:]
	mx, err := net.LookupMX(domain)
	if err != nil || len(mx) == 0 {
		return false
	}

	return true
}
