package service

import (
	"net/url"
	"os"
	"strings"

	"github.com/FoPQer/go-shortener/internal/config/flags"
)

var fileStoragePath string

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

func initFileStoragePath() string {
	var filePath string

	if path := os.Getenv("FILE_STORAGE_PATH"); path != "" {
		filePath = path
	} else {
		filePath = flags.GetFlagFileStoragePath()
	}
	
	file, err := os.Open(filePath)
	if err != nil {
		file, err = os.Create(filePath)
	}
	file.Close()

	fileStoragePath = filePath
	return filePath
}

func GetFileStoragePath() string {
	if fileStoragePath == "" {
		return initFileStoragePath()
	}
	return fileStoragePath
}
