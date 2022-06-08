package validation

func IsPasswordValid(password string, min int) bool {
	if len(password) < min {
		return false
	}
	return true
}
