package service

import (
	"flag"
	"io"
	"os"
	"testing"
)

func TestConfigPriority(t *testing.T) {
	// Reset cache before test
	resetConfigCache()

	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{
		"server_address": "config.example.com:9999",
		"base_url": "http://config.example.com",
		"file_storage_path": "/config/storage.json",
		"database_dsn": "postgres://config:config@localhost/db",
		"enable_https": true,
		"audit_file": "/config/audit.log",
		"audit_url": "http://config.example.com/audit"
	}`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Set environment variable for config file
	t.Setenv("CONFIG", tmpFile.Name())

	tests := []struct {
		name     string
		setFlag  func()
		setEnv   func()
		verify   func(t *testing.T)
		cleanup  func()
	}{
		{
			name: "config file only",
			setFlag: func() {
				// No flags set
			},
			setEnv: func() {
				// Config environment variable already set
			},
			verify: func(t *testing.T) {
				resetConfigCache()
				if addr := GetRunAddr(); addr != "config.example.com:9999" {
					t.Errorf("GetRunAddr() = %s, want config.example.com:9999", addr)
				}
			},
			cleanup: func() {
				resetConfigCache()
			},
		},
		{
			name: "environment variable overrides config",
			setFlag: func() {
				// No flags set
			},
			setEnv: func() {
				os.Setenv("SERVER_ADDRESS", "env.example.com:7777")
			},
			verify: func(t *testing.T) {
				resetConfigCache()
				if addr := GetRunAddr(); addr != "env.example.com:7777" {
					t.Errorf("GetRunAddr() = %s, want env.example.com:7777", addr)
				}
			},
			cleanup: func() {
				os.Unsetenv("SERVER_ADDRESS")
				resetConfigCache()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setFlag()
			tt.setEnv()
			tt.verify(t)
			tt.cleanup()
		})
	}
}

func TestGetBaseURLFromConfig(t *testing.T) {
	resetConfigCache()

	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{
		"base_url": "http://config.example.com"
	}`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	os.Setenv("CONFIG", tmpFile.Name())
	defer os.Unsetenv("CONFIG")

	baseURL := GetBasePrefix()
	if !contains(baseURL, "config.example.com") && !contains(baseURL, "example.com") {
		t.Errorf("GetBasePrefix() = %s, want config value", baseURL)
	}

	resetConfigCache()
}

func TestGetFileStoragePathFromConfig(t *testing.T) {
	resetConfigCache()

	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{
		"file_storage_path": "/config/storage.json"
	}`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	os.Setenv("CONFIG", tmpFile.Name())
	defer os.Unsetenv("CONFIG")

	storagePath := GetFileStoragePath()
	if storagePath != "/config/storage.json" {
		t.Errorf("GetFileStoragePath() = %s, want /config/storage.json", storagePath)
	}

	resetConfigCache()
}

func TestGetDatabaseDSNFromConfig(t *testing.T) {
	resetConfigCache()

	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{
		"database_dsn": "postgres://config:config@localhost/db"
	}`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	os.Setenv("CONFIG", tmpFile.Name())
	defer os.Unsetenv("CONFIG")

	dsn := GetDatabaseDSN()
	if dsn != "postgres://config:config@localhost/db" {
		t.Errorf("GetDatabaseDSN() = %s, want postgres://config:config@localhost/db", dsn)
	}

	resetConfigCache()
}

func TestGetHTTPSFromConfig(t *testing.T) {
	resetConfigCache()

	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{
		"enable_https": true
	}`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	os.Setenv("CONFIG", tmpFile.Name())
	defer os.Unsetenv("CONFIG")

	https := GetHTTPs()
	if !https {
		t.Errorf("GetHTTPs() = %v, want true", https)
	}

	resetConfigCache()
}

func TestIsFlagSetDoesNotMatchByPrefix(t *testing.T) {
	originalCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = originalCommandLine
	}()

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var runAddr string
	var auditFile string
	fs.StringVar(&runAddr, "a", "", "")
	fs.StringVar(&auditFile, "audit-file", "", "")

	if err := fs.Parse([]string{"-audit-file=/tmp/audit.log"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	flag.CommandLine = fs

	if isFlagSet("a") {
		t.Fatal("isFlagSet(\"a\") should be false when only -audit-file is set")
	}

	if !isFlagSet("audit-file") {
		t.Fatal("isFlagSet(\"audit-file\") should be true when -audit-file is set")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
