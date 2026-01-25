package handler

import (
	"io"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func GetURL(res http.ResponseWriter, req *http.Request) {
	urls := repository.Urls
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(res, "", 400)
		return
	}

	url, err := urls.GetURL(id)
	if err != nil {
		http.Error(res, "", 400)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(307)
}

func PostURL(res http.ResponseWriter, req *http.Request) {
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

	id := service.NewID()
	if err := urls.SetURL(id, string(body)); err != nil {
		http.Error(res, "", 400)
		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-Type", "text/plain")
	res.Write([]byte("http://" + flags.FlagRunAddr + flags.FlagBasePrefix + "/" + id))
}
