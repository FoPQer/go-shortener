package flags

import (
	"flag"
)

var (
	flagRunAddr         string
	flagBasePrefix      string
	flagFileStoragePath string
	flagDatabaseURL     string
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

func ParseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagBasePrefix, "b", "", "base prefix to URL")
	flag.StringVar(&flagFileStoragePath, "f", "", "file storage path")
	flag.StringVar(&flagDatabaseURL, "d", "", "database connection string")

	flag.Parse()
}
