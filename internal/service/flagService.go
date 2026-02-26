package service

import (
	"net/url"
	"os"
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
