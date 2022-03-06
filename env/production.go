package env

func IsProduction() bool {
	return GetEnvBoolDef("PRODUCTION", false)
}

func IsDevelopment() bool {
	return !IsProduction()
}

func EnvType() string {
	if IsProduction() {
		return "production"
	}

	return "development"
}
