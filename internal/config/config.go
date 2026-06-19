package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/FoPQer/go-shortener/internal/logger"
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
	TrustedSubnet    string `json:"trusted_subnet"`
}

// LoadConfig loads configuration from the specified file path.
// Returns an error if the file cannot be read or parsed.
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

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		if lg := logger.GetSugar(); lg != nil {
			lg.Errorf("Failed to read config file: %v", err)
		}
		return nil, err
	}

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		if lg := logger.GetSugar(); lg != nil {
			lg.Errorf("Failed to parse config file: %v", err)
		}
		return nil, err
	}

	if lg := logger.GetSugar(); lg != nil {
		lg.Infof("Configuration loaded from file: %s", filePath)
	}

	return &cfg, nil
}
