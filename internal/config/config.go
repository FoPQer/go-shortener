package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the configuration structure loaded from a JSON file.
type Config struct {
	ServerAddress    string `json:"server_address"`
	BaseURL          string `json:"base_url"`
	FileStoragePath  string `json:"file_storage_path"`
	DatabaseDSN      string `json:"database_dsn"`
	EnableHTTPS      bool   `json:"enable_https"`
	AuditFile        string `json:"audit_file"`
	AuditURL         string `json:"audit_url"`
}

var loadedConfig *Config

// LoadConfig loads configuration from the specified file path.
// Returns nil if the file doesn't exist or can't be read.
func LoadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		return nil, nil
	}

	// Normalize the file path
	filePath = strings.TrimSpace(filePath)
	filePath = strings.Trim(filePath, "\"'")
	if filePath == "" {
		return nil, nil
	}

	filePath = filepath.Clean(filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Config file not found: %s\n", filePath)
		return nil, nil
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config file: %v\n", err)
		return nil, err
	}

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config file: %v\n", err)
		return nil, err
	}

	loadedConfig = &cfg
	fmt.Fprintf(os.Stderr, "Configuration loaded from file: %s\n", filePath)

	return &cfg, nil
}

// GetLoadedConfig returns the currently loaded configuration or nil.
func GetLoadedConfig() *Config {
	return loadedConfig
}

// ResetConfig clears the loaded configuration.
func ResetConfig() {
	loadedConfig = nil
}
