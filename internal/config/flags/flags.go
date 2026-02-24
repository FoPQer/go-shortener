package flags

import (
	"flag"
)

var (
	flagRunAddr         string
	flagBasePrefix      string
	flagFileStoragePath string
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

func ParseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagBasePrefix, "b", "", "base prefix to URL")
	flag.StringVar(&flagFileStoragePath, "f", "settings.json", "file storage path")

	flag.Parse()
}
