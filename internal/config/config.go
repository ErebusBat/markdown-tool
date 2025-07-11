package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/erebusbat/markdown-tool/pkg/types"
	"github.com/spf13/viper"
)

// Load loads configuration from file or creates default config
func Load(configFile string) (*types.Config, error) {
	// Create a new Viper instance with custom key delimiter to handle domain names with dots
	// Using "::" instead of "." prevents domain names like "companycam.slack.com" 
	// from being interpreted as nested YAML structures
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// Set default config path
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		configDir := filepath.Join(home, ".config", "markdown-tool")
		configPath := filepath.Join(configDir, "config.yaml")

		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		// Create default config if it doesn't exist
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			if err := createDefaultConfig(configPath); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
		}

		v.SetConfigFile(configPath)
	}

	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config types.Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(path string) error {
	defaultConfig := `github:
  default_org: "CompanyCam"
  default_repo: "Company-Cam-API"
  mappings:
    "CompanyCam/Company-Cam-API": "CompanyCam/API"

jira:
  domain: "https://companycam.atlassian.net"
  projects:
    - "PLAT"
    - "SPEED"

url:
  domain_mappings:
    "companycam.slack.com": "slack"
    "youtube.com": "YouTube"
`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}
