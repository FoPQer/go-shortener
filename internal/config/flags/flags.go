package flags

import (
	"flag"
)

var (
	flagRunAddr         string
	flagBasePrefix      string
	flagFileStoragePath string
	flagDatabaseURL     string
	flagAuditFile       string
	flagAuditURL        string
	flagHTTPs           bool
	flagConfigFile      string
)

// GetFlagRunAddr returns server address configured via command-line flags.
func GetFlagRunAddr() string {
	return flagRunAddr
}

// SetFlagRunAddr sets server address flag value.
func SetFlagRunAddr(newFlagRunAddr string) {
	flagRunAddr = newFlagRunAddr
}

// GetFlagBasePrefix returns URL base prefix configured via command-line flags.
func GetFlagBasePrefix() string {
	return flagBasePrefix
}

// SetFlagBasePrefix sets URL base prefix flag value.
func SetFlagBasePrefix(newFlagBasePrefix string) {
	flagBasePrefix = newFlagBasePrefix
}

// GetFlagFileStoragePath returns file storage path configured via command-line flags.
func GetFlagFileStoragePath() string {
	return flagFileStoragePath
}

// SetFlagFileStoragePath sets file storage path flag value.
func SetFlagFileStoragePath(newFlagFileStoragePath string) {
	flagFileStoragePath = newFlagFileStoragePath
}

// GetFlagDatabaseURL returns database DSN configured via command-line flags.
func GetFlagDatabaseURL() string {
	return flagDatabaseURL
}

// SetFlagDatabaseURL sets database DSN flag value.
func SetFlagDatabaseURL(newFlagDatabaseURL string) {
	flagDatabaseURL = newFlagDatabaseURL
}

// GetFlagAuditFile returns audit file path configured via command-line flags.
func GetFlagAuditFile() string {
	return flagAuditFile
}

// SetFlagAuditFile sets audit file path flag value.
func SetFlagAuditFile(newFlagAuditFile string) {
	flagAuditFile = newFlagAuditFile
}

// GetFlagAuditURL returns audit endpoint URL configured via command-line flags.
func GetFlagAuditURL() string {
	return flagAuditURL
}

// SetFlagAuditURL sets audit endpoint URL flag value.
func SetFlagAuditURL(newFlagAuditURL string) {
	flagAuditURL = newFlagAuditURL
}

// GetFlagHTTPs returns HTTPS configuration configured via command-line flags.
func GetFlagHTTPs() bool {
	return flagHTTPs
}

// SetFlagHTTPs sets HTTPS configuration flag value.
func SetFlagHTTPs(newFlagHTTPs bool) {
	flagHTTPs = newFlagHTTPs
}

// GetFlagConfigFile returns config file path configured via command-line flags.
func GetFlagConfigFile() string {
	return flagConfigFile
}

// SetFlagConfigFile sets config file path flag value.
func SetFlagConfigFile(newFlagConfigFile string) {
	flagConfigFile = newFlagConfigFile
}

// ParseFlags defines and parses all supported command-line flags for the service.
func ParseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagBasePrefix, "b", "", "base prefix to URL")
	flag.StringVar(&flagFileStoragePath, "f", "", "file storage path")
	flag.StringVar(&flagDatabaseURL, "d", "", "database connection string")
	flag.StringVar(&flagAuditFile, "audit-file", "", "audit file path")
	flag.StringVar(&flagAuditURL, "audit-url", "", "audit URL")
	flag.BoolVar(&flagHTTPs, "s", false, "HTTPS configuration")
	flag.StringVar(&flagConfigFile, "c", "", "config file path")
	flag.Parse()
}
