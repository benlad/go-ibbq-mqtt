package main

import (
	"os"
	"strconv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvBool(key string) bool {
	logger.Debug("getEnvBool", "Looking up", key)
	if value, ok := os.LookupEnv(key); ok {
		logger.Debug("getEnvBool", "Key found, value", value)
		if boolValue, err := strconv.ParseBool(value); err == nil {
			logger.Debug("getEnvBool", "Bool value", boolValue)
			return boolValue
		}
	}
	return false
}
