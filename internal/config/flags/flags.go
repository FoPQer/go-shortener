package flags

import (
	"flag"
)

var (
	flagRunAddr         string
	flagBasePrefix      string
	flagFileStoragePath string
	flagDatabaseURL     string
	flagAuditFile		string
    flagAuditURL		string
)

func GetFlagRunAddr() string {
	return flagRunAddr
}

func SetFlagRunAddr(newFlagRunAddr string) {
	flagRunAddr = newFlagRunAddr
}

func GetFlagBasePrefix() string {
	return flagBasePrefix
}

func SetFlagBasePrefix(newFlagBasePrefix string) {
	flagBasePrefix = newFlagBasePrefix
}

func GetFlagFileStoragePath() string {
	return flagFileStoragePath
}

func SetFlagFileStoragePath(newFlagFileStoragePath string) {
	flagFileStoragePath = newFlagFileStoragePath
}

func GetFlagDatabaseURL() string {
	return flagDatabaseURL
}

func SetFlagDatabaseURL(newFlagDatabaseURL string) {
	flagDatabaseURL = newFlagDatabaseURL
}

func GetFlagAuditFile() string {
	return flagAuditFile
}

func SetFlagAuditFile(newFlagAuditFile string) {
	flagAuditFile = newFlagAuditFile
}

func GetFlagAuditURL() string {
	return flagAuditURL
}

func SetFlagAuditURL(newFlagAuditURL string) {
	flagAuditURL = newFlagAuditURL
}

func ParseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagBasePrefix, "b", "", "base prefix to URL")
	flag.StringVar(&flagFileStoragePath, "f", "", "file storage path")
	flag.StringVar(&flagDatabaseURL, "d", "", "database connection string")
	flag.StringVar(&flagAuditFile, "audit-file", "", "audit file path")
	flag.StringVar(&flagAuditURL, "audit-url", "", "audit URL")

	flag.Parse()
}
