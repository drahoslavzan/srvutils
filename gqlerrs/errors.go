package gqlerrs

func UnauthorizedError() SrvError {
	return MakeSrvError("auth/unauthorized", nil)
}
