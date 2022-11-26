package env

func IsProduction() bool {
	return !IsDevelopment()
}

func IsDevelopment() bool {
	return BoolDef("ENV_DEVELOPMENT", false)
}

func Type() string {
	if IsProduction() {
		return "production"
	}

	return "development"
}
