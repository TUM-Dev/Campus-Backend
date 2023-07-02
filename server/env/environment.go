package env

import "os"

const (
	DefaultApiKey = "default-api-key"
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

func ApiKey() string {
	value, ok := os.LookupEnv("API_KEY")

	if !ok {
		return DefaultApiKey
	}

	return value
}

func ApiUrl() string {
	value, ok := os.LookupEnv("API_URL")

	if !ok {
		return "http://localhost:50051"
	}

	return value
}
