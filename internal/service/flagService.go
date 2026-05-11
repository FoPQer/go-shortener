package service

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/FoPQer/go-shortener/internal/config/flags"
)

func GetRunAddr() string {
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		return url.PathEscape(addr)
	} else {
		return url.PathEscape(flags.GetFlagRunAddr())
	}
}

func GetBasePrefix() string {
	var outBase string

	if base := os.Getenv("BASE_URL"); base != "" {
		outBase = base
	} else {
		outBase = flags.GetFlagBasePrefix()
	}

	outBase = url.PathEscape(outBase)

	if !strings.HasPrefix(outBase, "/") {
		outBase = "/" + outBase
	}
	if !strings.HasSuffix(outBase, "/") {
		outBase = outBase + "/"
	}

	return outBase
}

func GetFileStoragePath() string {
	if path := os.Getenv("FILE_STORAGE_PATH"); path != "" {
		return path
	} else {
		return flags.GetFlagFileStoragePath()
	}
}

func GetDatabaseDSN() string {
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		return dsn
	} else {
		return flags.GetFlagDatabaseURL()
	}
}

func GetSecretKey() string {
	if secretKey := os.Getenv("SECRET_KEY"); secretKey != "" {
		return secretKey
	} else {
		return "your_secret_key";
	}
}

func GetAuditFile() string {
	if auditFile := os.Getenv("AUDIT_FILE"); auditFile != "" {
		return normalizePath(auditFile)
	} else {
		return normalizePath(flags.GetFlagAuditFile())
	}
}

func GetAuditURL() string {
	if auditURL := os.Getenv("AUDIT_URL"); auditURL != "" {
		return auditURL
	} else {
		return flags.GetFlagAuditURL()
	}
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "\"'")
	if path == "" {
		return ""
	}

	return filepath.Clean(path)
}
