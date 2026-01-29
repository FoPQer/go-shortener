package flags

import (
	"flag"
	"strings"
)

var (
	FlagRunAddr    string
	FlagBasePrefix string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBasePrefix, "b", "", "base prefix to URL")

	flag.Parse()
	if !strings.HasPrefix(FlagBasePrefix, "/") {
		FlagBasePrefix = "/" + FlagBasePrefix
	}
}
