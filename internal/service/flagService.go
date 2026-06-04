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

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "\"'")
	if path == "" {
		return ""
	}

	return filepath.Clean(path)
}
