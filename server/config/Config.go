package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds the server configuration.
type Config struct {
	ReadQuorum  int
	WriteQuorum int
	Redundancy  int
}

func NewConfig() *Config {
	return &Config{
		ReadQuorum:  getEnvAsInt("READ_QUORUM", 2),
		WriteQuorum: getEnvAsInt("WRITE_QUORUM", 2),
		Redundancy:  getEnvAsInt("REDUNDANCY", 3),
	}
}

// getEnvAsInt reads an environment variable and converts it to an integer.
// If the variable is not set or cannot be converted, it returns the default value.
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s: %v, using default %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}
