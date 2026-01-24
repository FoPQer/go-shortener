package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/FoPQer/go-shortener/internal/service"
)

func GetUrl(res http.ResponseWriter, req *http.Request) {
	urls := repository.Urls
	splittedPath := strings.Split(strings.TrimPrefix(req.URL.Path, "/"), "/")
	if len(splittedPath) > 1 {
		http.Error(res, "", 400)
		return
	}

	url, err := urls.GetUrl(splittedPath[0])
	if err != nil {
		http.Error(res, "", 400)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(307)
}

func PostUrl(res http.ResponseWriter, req *http.Request) {
	urls := repository.Urls
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "", 400)
		return
	}
	if len(body) == 0 {
		http.Error(res, "", 400)
		return
	}

	id := service.NewId()
	if err := urls.SetUrl(id, string(body)); err != nil {
		http.Error(res, "", 400)
		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-Type", "text/plain")
	res.Write([]byte("http://localhost:8080/" + id))
}
