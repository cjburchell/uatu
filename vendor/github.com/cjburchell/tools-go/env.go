package tools

import (
	"os"
	"strconv"
)

func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if result, err := strconv.Atoi(value); err == nil {
			return result
		}
		return fallback
	}
	return fallback
}

func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if result, err := strconv.ParseBool(value); err == nil {
			return result
		}
	}
	return fallback
}
