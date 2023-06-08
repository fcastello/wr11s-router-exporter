// config/config.go

package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Config holds the application configuration
type Config struct {
	Address  string
	Username string
	Password string
}

// LoadConfig loads the configuration from CLI flags, environment variables, and TOML file
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Initialize logger
	logger := logrus.New()

	// Load configuration from TOML file
	err := loadConfigFromToml(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load TOML configuration: %w", err)
	}

	// Load configuration from environment variables
	err = loadConfigFromEnv(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load environment variable configuration: %w", err)
	}

	// Load configuration from CLI flags
	err = loadConfigFromFlags(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load CLI flag configuration: %w", err)
	}

	// Validate required fields
	err = validateConfig(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return config, nil
}

// Load configuration from TOML file
func loadConfigFromToml(config *Config, logger *logrus.Logger) error {
	// Load configuration file (assuming it's named "config.toml")
	file, err := os.Open("config.toml")
	if err != nil {
		logger.WithError(err).Error("failed to open TOML configuration file")
		return err
	}
	defer file.Close()

	// Parse the TOML data
	tree, err := toml.LoadReader(file)
	if err != nil {
		logger.WithError(err).Error("failed to parse TOML configuration")
		return err
	}

	// Get the values from TOML
	config.Address = tree.Get("address").(string)
	config.Username = tree.Get("username").(string)
	config.Password = tree.Get("password").(string)

	logger.WithFields(logrus.Fields{
		"address":  config.Address,
		"username": config.Username,
		"password": "*****", // Mask password in logs
	}).Debug("Loaded configuration from TOML")

	return nil
}

// Load configuration from environment variables
func loadConfigFromEnv(config *Config, logger *logrus.Logger) error {
	// Get environment variables
	address := os.Getenv("WR11_ADDRESS")
	username := os.Getenv("WR11_USERNAME")
	password := os.Getenv("WR11_PASSWORD")

	// Override config values with environment variables if they exist
	if address != "" {
		config.Address = address
	}
	if username != "" {
		config.Username = username
	}
	if password != "" {
		config.Password = password
	}

	logger.WithFields(logrus.Fields{
		"address":  config.Address,
		"username": config.Username,
		"password": "*****", // Mask password in logs
	}).Debug("Loaded configuration from environment variables")

	return nil
}

// Load configuration from CLI flags
func loadConfigFromFlags(config *Config, logger *logrus.Logger) error {
	// Define and parse CLI flags
	flagSet := pflag.NewFlagSet("app", pflag.ExitOnError)
	flagSet.StringVar(&config.Address, "address", "", "Base URL for HTTP requests")
	flagSet.StringVar(&config.Username, "username", "", "Login username")
	flagSet.StringVar(&config.Password, "password", "", "Login password")
	flagSet.Parse(os.Args[1:])

	logger.WithFields(logrus.Fields{
		"address":  config.Address,
		"username": config.Username,
		"password": "*****", // Mask password in logs
	}).Debug("Loaded configuration from CLI flags")

	return nil
}

// Validate required fields
func validateConfig(config *Config, logger *logrus.Logger) error {
	if config.Address == "" {
		logger.Error("address is required")
		return fmt.Errorf("address is required")
	}
	if config.Username == "" {
		logger.Error("username is required")
		return fmt.Errorf("username is required")
	}
	if config.Password == "" {
		logger.Error("password is required")
		return fmt.Errorf("password is required")
	}
	return nil
}
