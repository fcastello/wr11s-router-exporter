// config/config_test.go

package config

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromToml(t *testing.T) {
	// Set up a temporary TOML file for testing
	tomlData := `
		address = "http://example.com"
		username = "testuser"
		password = "testpassword"
	`
	tmpFile, err := os.CreateTemp("", "config*.toml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(tomlData)
	assert.NoError(t, err)

	// Set the environment variable for TOML file path
	os.Setenv("WR11_ADDRESS", "http://envaddress.com")
	os.Setenv("WR11_USERNAME", "envuser")
	os.Setenv("WR11_PASSWORD", "envpassword")

	// Initialize logger
	logger := logrus.New()

	// Load configuration
	config, err := LoadConfig(logger)
	assert.NoError(t, err)

	// Validate the values loaded from TOML have lower precedence
	assert.Equal(t, "http://example.com", config.Address)
	assert.Equal(t, "testuser", config.Username)
	assert.Equal(t, "testpassword", config.Password)
}

// Similar tests can be written for loading configuration from environment variables and CLI flags.

func TestMissingVariables(t *testing.T) {
	// Set the environment variables for testing
	os.Setenv("WR11_ADDRESS", "")
	os.Setenv("WR11_USERNAME", "")
	os.Setenv("WR11_PASSWORD", "")

	// Initialize logger
	logger := logrus.New()

	// Load configuration with missing variables
	_, err := LoadConfig(logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "address is required")

	// Reset environment variables
	os.Setenv("WR11_ADDRESS", "http://example.com")
	os.Setenv("WR11_USERNAME", "")
	os.Setenv("WR11_PASSWORD", "")

	// Load configuration with missing variables
	_, err = LoadConfig(logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username is required")

	// Reset environment variables
	os.Setenv("WR11_ADDRESS", "http://example.com")
	os.Setenv("WR11_USERNAME", "testuser")
	os.Setenv("WR11_PASSWORD", "")

	// Load configuration with missing variables
	_, err = LoadConfig(logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password is required")
}
