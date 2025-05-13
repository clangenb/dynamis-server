package utils

import (
	"os"
	"strings"
)

const AppEnvEnvVar = "APP_ENV"

const AppEnvDev = "dev"
const AppEnvProd = "prod"

func GetAppEnv() string {
	env := strings.ToLower(os.Getenv(AppEnvEnvVar))
	if env == "" {
		return AppEnvDev // sensible default
	}
	return env
}

func IsDev() bool {
	return GetAppEnv() == AppEnvDev
}

func IsProd() bool {
	return GetAppEnv() == AppEnvProd
}
