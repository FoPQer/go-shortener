package service

import (
	"os"
	"strings"

	"github.com/FoPQer/go-shortener/internal/config/flags"
)

func GetRunAddr() string {
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		return addr
	} else {
		return flags.GetFlagRunAddr()
	}
}

func GetBasePrefix() string {
	var out_base string

	if base := os.Getenv("BASE_URL"); base != "" {
		out_base = base
	} else {
		out_base = flags.GetFlagBasePrefix()
	}

	if !strings.HasPrefix(out_base, "/") {
		out_base = "/" + out_base
	}
	if !strings.HasSuffix(out_base, "/") {
		out_base = out_base + "/"
	}

	return out_base
}
