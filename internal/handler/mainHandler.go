package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func GetURL(res http.ResponseWriter, req *http.Request) {
	urls := repository.GetUrls()
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
	urls := repository.GetUrls()
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

	target, err := url.JoinPath("http://"+flags.GetFlagRunAddr(), flags.GetFlagBasePrefix(), id)

	if err != nil {
		http.Error(res, "", 400)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(target))
}
