package env

import "os"

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
