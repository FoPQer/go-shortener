package service

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/FoPQer/go-shortener/internal/config/flags"
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
)

// GetRunAddr returns the configured server address.
func GetRunAddr() string {
	runAddrOnce.Do(func() {
		if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
			runAddr = url.PathEscape(addr)
			return
		}

		runAddr = url.PathEscape(flags.GetFlagRunAddr())
	})

	return runAddr
}

// GetBasePrefix returns the configured URL base prefix.
func GetBasePrefix() string {
	basePrefixOnce.Do(func() {
		if base := os.Getenv("BASE_URL"); base != "" {
			basePrefix = base
		} else {
			basePrefix = flags.GetFlagBasePrefix()
		}

		basePrefix = url.PathEscape(basePrefix)

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
		if path := os.Getenv("FILE_STORAGE_PATH"); path != "" {
			fileStoragePath = path
			return
		}

		fileStoragePath = flags.GetFlagFileStoragePath()
	})

	return fileStoragePath
}

// GetDatabaseDSN returns the configured database DSN.
func GetDatabaseDSN() string {
	databaseDSNOnce.Do(func() {
		if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
			databaseDSN = dsn
			return
		}

		databaseDSN = flags.GetFlagDatabaseURL()
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
		if value := os.Getenv("AUDIT_FILE"); value != "" {
			auditFile = normalizePath(value)
			return
		}

		auditFile = normalizePath(flags.GetFlagAuditFile())
	})

	return auditFile
}

// GetAuditURL returns the configured audit endpoint URL.
func GetAuditURL() string {
	auditURLOnce.Do(func() {
		if value := os.Getenv("AUDIT_URL"); value != "" {
			auditURL = value
			return
		}

		auditURL = flags.GetFlagAuditURL()
	})

	return auditURL
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
