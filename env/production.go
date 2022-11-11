package env

func IsProduction() bool {
	return !IsDevelopment()
}

func IsDevelopment() bool {
	return GetEnvBoolDef("ENV_DEVELOPMENT", false)
}

func EnvType() string {
	if IsProduction() {
		return "production"
	}

	return "development"
}
