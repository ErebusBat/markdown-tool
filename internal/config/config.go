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
	if configFile != "" {
		viper.SetConfigFile(configFile)
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
		
		viper.SetConfigFile(configPath)
	}

	viper.SetConfigType("yaml")
	
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
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
`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}