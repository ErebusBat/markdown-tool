package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "markdown-tool")
	configPath := filepath.Join(configDir, "config.yaml")

	// Mock home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Load config (should create default)
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify config directory was created
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Config directory was not created")
	}

	// Verify config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Verify default values
	if cfg.GitHub.DefaultOrg != "CompanyCam" {
		t.Errorf("GitHub.DefaultOrg = %q, want CompanyCam", cfg.GitHub.DefaultOrg)
	}

	if cfg.GitHub.DefaultRepo != "Company-Cam-API" {
		t.Errorf("GitHub.DefaultRepo = %v, want Company-Cam-API", cfg.GitHub.DefaultRepo)
	}

	if cfg.JIRA.Domain != "https://companycam.atlassian.net" {
		t.Errorf("JIRA.Domain = %v, want https://companycam.atlassian.net", cfg.JIRA.Domain)
	}

	expectedProjects := []string{"PLAT", "SPEED"}
	if len(cfg.JIRA.Projects) != len(expectedProjects) {
		t.Errorf("JIRA.Projects length = %v, want %v", len(cfg.JIRA.Projects), len(expectedProjects))
	}

	for i, project := range expectedProjects {
		if i >= len(cfg.JIRA.Projects) || cfg.JIRA.Projects[i] != project {
			t.Errorf("JIRA.Projects[%v] = %v, want %v", i, cfg.JIRA.Projects[i], project)
		}
	}

	// Verify GitHub mappings (note: Viper lowercases map keys)
	if len(cfg.GitHub.Mappings) == 0 {
		t.Error("GitHub.Mappings should not be empty")
	}
}

func TestLoad_CustomConfigFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "custom-config.yaml")

	customConfig := `github:
  default_org: "TestOrg"
  default_repo: "TestRepo"
  mappings:
    "testorg/testrepo": "Test/Repo"

jira:
  domain: "https://test.atlassian.net"
  projects:
    - "TEST"
`

	err := os.WriteFile(configPath, []byte(customConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load custom config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify custom values
	if cfg.GitHub.DefaultOrg != "TestOrg" {
		t.Errorf("GitHub.DefaultOrg = %v, want TestOrg", cfg.GitHub.DefaultOrg)
	}

	if cfg.GitHub.DefaultRepo != "TestRepo" {
		t.Errorf("GitHub.DefaultRepo = %v, want TestRepo", cfg.GitHub.DefaultRepo)
	}

	if cfg.JIRA.Domain != "https://test.atlassian.net" {
		t.Errorf("JIRA.Domain = %v, want https://test.atlassian.net", cfg.JIRA.Domain)
	}

	if len(cfg.JIRA.Projects) != 1 || cfg.JIRA.Projects[0] != "TEST" {
		t.Errorf("JIRA.Projects = %v, want [TEST]", cfg.JIRA.Projects)
	}
}

func TestLoad_InvalidConfigFile(t *testing.T) {
	// Create a temporary invalid config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid-config.yaml")

	invalidConfig := `invalid: yaml: content: [unclosed`

	err := os.WriteFile(configPath, []byte(invalidConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load should fail with invalid config
	_, err = Load(configPath)
	if err == nil {
		t.Error("Expected error loading invalid config file")
	}
}
