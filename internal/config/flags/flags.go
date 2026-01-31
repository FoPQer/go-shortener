package flags

import (
	"flag"
	"net/url"
)

var (
	flagRunAddr    string
	flagBasePrefix string
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

func ParseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagBasePrefix, "b", "", "base prefix to URL")

	flag.Parse()

	flagBasePrefix = url.PathEscape(flagBasePrefix)
}
