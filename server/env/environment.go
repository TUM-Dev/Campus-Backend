package env

import (
	"github.com/guregu/null"
	"os"
)

func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

func IsDev() bool {
	return GetEnvironment() == "dev"
}

func IsProd() bool {
	return GetEnvironment() == "prod"
}

func IsMensaCronActive() bool {
	return os.Getenv("MensaCronDisabled") != "true"
}

func ApiKey() null.String {
	return null.StringFrom(os.Getenv("API_KEY"))
}
