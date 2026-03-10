package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/FoPQer/go-shortener/internal/config/db"
	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	urlService *service.URLService
	jsonService *service.JSONService
}

func NewHandler(urlService *service.URLService, jsonService *service.JSONService) *Handler {
	return &Handler{urlService: urlService, jsonService: jsonService}
}

func (h *Handler) GetPing(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    if err := db.GetDBConn().Ping(ctx); err != nil {
        http.Error(res, "", http.StatusInternalServerError)
        return
    }
	res.WriteHeader(http.StatusOK)
}

func (h *Handler) GetURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	logger.GetSugar().Infof("%v", h.urlService.GetUrls())

	url, err := h.urlService.GetURL(id)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	logger.GetSugar().Infof("%s", url)
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) PostURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	logger.GetSugar().Infof("body: %s", string(body))
	target, err := h.urlService.SetURL(string(body))
	if errors.Is(err, model.ErrURLAlreadyExists) {
		res.WriteHeader(http.StatusConflict)
	} else if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	logger.GetSugar().Infof("target: %s", target)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(target))
}

func (h *Handler) PostURLByJSON(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	url, err := h.jsonService.GetURLFromJSON(body)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	target, err := h.urlService.SetURL(string(url))
	if errors.Is(err, model.ErrURLAlreadyExists) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusConflict)
	} else if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	out, err := h.jsonService.SetURLToJSON(target)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(out)
}

func (h *Handler) PostBatchURLByJSON(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	urls, err := h.jsonService.GetBatchURLFromJSON(body)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	targets, err := h.urlService.SetBatchURL(urls)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	out, err := h.jsonService.SetBatchURLToJSON(targets)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(out)
}
