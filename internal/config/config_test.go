package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		filePath  string
		wantError bool
		wantCfg   *Config
	}{
		{
			name: "valid config file",
			content: `{
				"server_address": "localhost:9090",
				"base_url": "http://example.com",
				"file_storage_path": "/tmp/storage",
				"database_dsn": "postgres://user:pass@localhost/db",
				"enable_https": true,
				"audit_file": "/var/log/audit.log",
				"audit_url": "http://audit.example.com"
			}`,
			filePath:  "test_config.json",
			wantError: false,
			wantCfg: &Config{
				ServerAddress:   "localhost:9090",
				BaseURL:         "http://example.com",
				FileStoragePath: "/tmp/storage",
				DatabaseDSN:     "postgres://user:pass@localhost/db",
				EnableHTTPS:     true,
				AuditFile:       "/var/log/audit.log",
				AuditURL:        "http://audit.example.com",
			},
		},
		{
			name:      "nonexistent file",
			filePath:  "/nonexistent/path/config.json",
			wantError: true,
		},
		{
			name:      "empty file path",
			filePath:  "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file if needed
			if tt.content != "" {
				tmpFile, err := os.CreateTemp("", "*.json")
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
				defer os.Remove(tmpFile.Name())

				if _, err := tmpFile.WriteString(tt.content); err != nil {
					t.Fatalf("failed to write to temp file: %v", err)
				}
				tmpFile.Close()
				tt.filePath = tmpFile.Name()
			}

			cfg, err := LoadConfig(tt.filePath)

			if (err != nil) != tt.wantError {
				t.Errorf("LoadConfig() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantCfg != nil && cfg != nil {
				if cfg.ServerAddress != tt.wantCfg.ServerAddress {
					t.Errorf("ServerAddress = %v, want %v", cfg.ServerAddress, tt.wantCfg.ServerAddress)
				}
				if cfg.BaseURL != tt.wantCfg.BaseURL {
					t.Errorf("BaseURL = %v, want %v", cfg.BaseURL, tt.wantCfg.BaseURL)
				}
				if cfg.FileStoragePath != tt.wantCfg.FileStoragePath {
					t.Errorf("FileStoragePath = %v, want %v", cfg.FileStoragePath, tt.wantCfg.FileStoragePath)
				}
				if cfg.DatabaseDSN != tt.wantCfg.DatabaseDSN {
					t.Errorf("DatabaseDSN = %v, want %v", cfg.DatabaseDSN, tt.wantCfg.DatabaseDSN)
				}
				if cfg.EnableHTTPS != tt.wantCfg.EnableHTTPS {
					t.Errorf("EnableHTTPS = %v, want %v", cfg.EnableHTTPS, tt.wantCfg.EnableHTTPS)
				}
				if cfg.AuditFile != tt.wantCfg.AuditFile {
					t.Errorf("AuditFile = %v, want %v", cfg.AuditFile, tt.wantCfg.AuditFile)
				}
				if cfg.AuditURL != tt.wantCfg.AuditURL {
					t.Errorf("AuditURL = %v, want %v", cfg.AuditURL, tt.wantCfg.AuditURL)
				}
			}
		})
	}
}

func TestLoadConfigIsStateless(t *testing.T) {
	// Create a temp config file
	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{"server_address": "localhost:8080"}`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	first, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("first LoadConfig() failed: %v", err)
	}

	second, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("second LoadConfig() failed: %v", err)
	}

	if first == nil || second == nil {
		t.Fatal("LoadConfig() should return non-nil config")
	}

	if first == second {
		t.Error("LoadConfig() should not rely on shared package-level state")
	}
}
