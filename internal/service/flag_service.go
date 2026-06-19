package service

import (
	"flag"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/FoPQer/go-shortener/internal/config"
	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/logger"
)

var (
	runAddrOnce sync.Once
	runAddr     string

	basePrefixOnce sync.Once
	basePrefix     string

	fileStoragePathOnce sync.Once
	fileStoragePath     string

	databaseDSNOnce sync.Once
	databaseDSN     string

	secretKeyOnce sync.Once
	secretKey     string

	auditFileOnce sync.Once
	auditFile     string

	auditURLOnce sync.Once
	auditURL     string

	httpsOnce sync.Once
	https     bool

	configOnce sync.Once
	cfg        *config.Config
)

// loadConfig loads configuration from file path.
// It uses sync.Once to ensure it's loaded only once.
func loadConfig() *config.Config {
	configOnce.Do(func() {
		configFilePath := getConfigFilePath()
		if configFilePath == "" {
			return
		}

		var err error
		cfg, err = config.LoadConfig(configFilePath)
		if err != nil {
			logger.GetSugar().Warnf("Failed to load config file: %v", err)
		}
	})

	return cfg
}

// getConfigFilePath returns the config file path from flags or environment variable.
func getConfigFilePath() string {
	// Check environment variable first
	if envPath := os.Getenv("CONFIG"); envPath != "" {
		return envPath
	}

	// Check command-line flag
	return flags.GetFlagConfigFile()
}

// isFlagSet checks if a flag was explicitly set in os.Args.
func isFlagSet(flagName string) bool {
	isSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			isSet = true
		}
	})

	return isSet
}

// GetRunAddr returns the configured server address.
func GetRunAddr() string {
	runAddrOnce.Do(func() {
		// Priority 1: Command-line flag
		if isFlagSet("a") {
			runAddr = url.PathEscape(flags.GetFlagRunAddr())
			return
		}

		// Priority 2: Environment variable
		if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
			runAddr = url.PathEscape(addr)
			return
		}

		// Priority 3: Config file
		cfg := loadConfig()
		if cfg != nil && cfg.ServerAddress != "" {
			runAddr = url.PathEscape(cfg.ServerAddress)
			return
		}

		// Priority 4: Default from flag
		runAddr = url.PathEscape(flags.GetFlagRunAddr())
	})

	return runAddr
}

// GetBasePrefix returns the configured URL base prefix.
func GetBasePrefix() string {
	basePrefixOnce.Do(func() {
		var baseValue string

		// Priority 1: Command-line flag
		if isFlagSet("b") {
			baseValue = flags.GetFlagBasePrefix()
		} else if base := os.Getenv("BASE_URL"); base != "" {
			// Priority 2: Environment variable
			baseValue = base
		} else {
			// Priority 3: Config file
			cfg := loadConfig()
			if cfg != nil && cfg.BaseURL != "" {
				baseValue = cfg.BaseURL
			} else {
				// Priority 4: Flag default
				baseValue = flags.GetFlagBasePrefix()
			}
		}

		basePrefix = url.PathEscape(baseValue)

		if !strings.HasPrefix(basePrefix, "/") {
			basePrefix = "/" + basePrefix
		}
		if !strings.HasSuffix(basePrefix, "/") {
			basePrefix = basePrefix + "/"
		}
	})

	return basePrefix
}

// GetFileStoragePath returns the configured file storage path.
func GetFileStoragePath() string {
	fileStoragePathOnce.Do(func() {
		var pathValue string

		// Priority 1: Command-line flag
		if isFlagSet("f") {
			pathValue = flags.GetFlagFileStoragePath()
		} else if path := os.Getenv("FILE_STORAGE_PATH"); path != "" {
			// Priority 2: Environment variable
			pathValue = path
		} else {
			// Priority 3: Config file
			cfg := loadConfig()
			if cfg != nil && cfg.FileStoragePath != "" {
				pathValue = cfg.FileStoragePath
			} else {
				// Priority 4: Flag default
				pathValue = flags.GetFlagFileStoragePath()
			}
		}

		fileStoragePath = pathValue
	})

	return fileStoragePath
}

// GetDatabaseDSN returns the configured database DSN.
func GetDatabaseDSN() string {
	databaseDSNOnce.Do(func() {
		var dsnValue string

		// Priority 1: Command-line flag
		if isFlagSet("d") {
			dsnValue = flags.GetFlagDatabaseURL()
		} else if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
			// Priority 2: Environment variable
			dsnValue = dsn
		} else {
			// Priority 3: Config file
			cfg := loadConfig()
			if cfg != nil && cfg.DatabaseDSN != "" {
				dsnValue = cfg.DatabaseDSN
			} else {
				// Priority 4: Flag default
				dsnValue = flags.GetFlagDatabaseURL()
			}
		}

		databaseDSN = dsnValue
	})

	return databaseDSN
}

// GetSecretKey returns the JWT secret key.
func GetSecretKey() string {
	secretKeyOnce.Do(func() {
		if value := os.Getenv("SECRET_KEY"); value != "" {
			secretKey = value
			return
		}

		secretKey = "your_secret_key"
	})

	return secretKey
}

// GetAuditFile returns the configured audit file path.
func GetAuditFile() string {
	auditFileOnce.Do(func() {
		var fileValue string

		// Priority 1: Command-line flag
		if isFlagSet("audit-file") {
			fileValue = flags.GetFlagAuditFile()
		} else if value := os.Getenv("AUDIT_FILE"); value != "" {
			// Priority 2: Environment variable
			fileValue = value
		} else {
			// Priority 3: Config file
			cfg := loadConfig()
			if cfg != nil && cfg.AuditFile != "" {
				fileValue = cfg.AuditFile
			} else {
				// Priority 4: Flag default
				fileValue = flags.GetFlagAuditFile()
			}
		}

		auditFile = normalizePath(fileValue)
	})

	return auditFile
}

// GetAuditURL returns the configured audit endpoint URL.
func GetAuditURL() string {
	auditURLOnce.Do(func() {
		var urlValue string

		// Priority 1: Command-line flag
		if isFlagSet("audit-url") {
			urlValue = flags.GetFlagAuditURL()
		} else if value := os.Getenv("AUDIT_URL"); value != "" {
			// Priority 2: Environment variable
			urlValue = value
		} else {
			// Priority 3: Config file
			cfg := loadConfig()
			if cfg != nil && cfg.AuditURL != "" {
				urlValue = cfg.AuditURL
			} else {
				// Priority 4: Flag default
				urlValue = flags.GetFlagAuditURL()
			}
		}

		auditURL = urlValue
	})

	return auditURL
}

// GetHTTPs returns the configured HTTPS setting.
func GetHTTPs() bool {
	httpsOnce.Do(func() {
		// Priority 1: Command-line flag
		if isFlagSet("s") {
			https = flags.GetFlagHTTPs()
			return
		}

		// Priority 2: Environment variable
		if value := os.Getenv("ENABLE_HTTPS"); value != "" {
			https = value == "true"
			return
		}

		// Priority 3: Config file
		cfg := loadConfig()
		if cfg != nil {
			https = cfg.EnableHTTPS
			return
		}

		// Priority 4: Flag default
		https = flags.GetFlagHTTPs()
	})

	return https
}

// resetConfigCache clears cached configuration values.
func resetConfigCache() {
	runAddrOnce = sync.Once{}
	runAddr = ""

	basePrefixOnce = sync.Once{}
	basePrefix = ""

	fileStoragePathOnce = sync.Once{}
	fileStoragePath = ""

	databaseDSNOnce = sync.Once{}
	databaseDSN = ""

	secretKeyOnce = sync.Once{}
	secretKey = ""

	auditFileOnce = sync.Once{}
	auditFile = ""

	auditURLOnce = sync.Once{}
	auditURL = ""

	httpsOnce = sync.Once{}
	https = false

	configOnce = sync.Once{}
	cfg = nil
}

// normalizePath trims quotes and whitespace and returns a cleaned file path.
func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "\"'")
	if path == "" {
		return ""
	}

	return filepath.Clean(path)
}
