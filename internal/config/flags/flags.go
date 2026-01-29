package flags

import (
	"flag"
	"net/url"
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
	FlagBasePrefix = url.QueryEscape(FlagBasePrefix)
	if !strings.HasPrefix(FlagBasePrefix, "/") {
		FlagBasePrefix = "/" + FlagBasePrefix
	}
}
