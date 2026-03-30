package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type OutputUserUrlsJson struct {
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Handler struct {
	urlService *service.URLService
	jsonService *service.JSONService
	userService *service.UserService
}

func NewHandler(urlService *service.URLService, jsonService *service.JSONService, userService *service.UserService) *Handler {
	return &Handler{urlService: urlService, jsonService: jsonService, userService: userService}
}

func (h *Handler) GetURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		logger.GetSugar().Errorln("Error while getting by shortUrl: empty id")
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	url, err := h.urlService.GetURL(id)
	if err != nil {
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) PostURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.GetSugar().Errorf("Error while reading request body: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		logger.GetSugar().Errorln("Error while reading request body: empty body")
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	logger.GetSugar().Infof("body: %s", string(body))
	userID := getUserIDFromContext(req.Context())
	target, err := h.urlService.SetURL(string(body), userID)
	if errors.Is(err, urls.ErrURLAlreadyExists) {
		res.WriteHeader(http.StatusConflict)
	} else if err != nil {
		logger.GetSugar().Errorf("Error while setting URL: %w", err)
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
		logger.GetSugar().Errorf("Error while reading request body: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		logger.GetSugar().Errorln("Error while reading request body: empty body")
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	url, err := h.jsonService.GetURLFromJSON(body)
	if err != nil {
		logger.GetSugar().Errorf("Error while getting URL from JSON: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	userID := getUserIDFromContext(req.Context())
	target, err := h.urlService.SetURL(string(url), userID)
	if errors.Is(err, urls.ErrURLAlreadyExists) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusConflict)
	} else if err != nil {
		logger.GetSugar().Errorf("Error while setting URL: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	out, err := h.jsonService.SetURLToJSON(target)
	if err != nil {
		logger.GetSugar().Errorf("Error while setting URL to JSON: %w", err)
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
		logger.GetSugar().Errorf("Error while reading request body: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		logger.GetSugar().Errorln("Error while reading request body: empty body")
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	urls, err := h.jsonService.GetBatchURLFromJSON(body)
	if err != nil {
		logger.GetSugar().Errorf("Error while getting batch url from JSON: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	userID := getUserIDFromContext(req.Context())
	targets, err := h.urlService.SetBatchURL(urls, userID)
	if err != nil {
		logger.GetSugar().Errorf("Error while setting batch url: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	out, err := h.jsonService.SetBatchURLToJSON(targets)
	if err != nil {
		logger.GetSugar().Errorf("Error while setting batch url to JSON: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(out)
}

func (h *Handler) GetUserURLs(res http.ResponseWriter, req *http.Request) {
	userID := getUserIDFromContext(req.Context())
	if userID == "" {
		http.Error(res, "Missing user ID", http.StatusUnauthorized)
		return
	}
	logger.GetSugar().Infof("UserID: %s", userID)
	urls, err := h.urlService.GetUrlsByUserID(userID)
	if err != nil {
		logger.GetSugar().Errorf("Error while getting user URLs: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	out, err := setUserUrlsToJSON(urls)
	if err != nil {
		logger.GetSugar().Errorf("Error while setting user URLs to JSON: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(out)
}

func (h *Handler) DeleteUserURLs(res http.ResponseWriter, req *http.Request) {
	userID := getUserIDFromContext(req.Context())
	if userID == "" {
		http.Error(res, "Missing user ID", http.StatusUnauthorized)
		return
	}
	logger.GetSugar().Infof("UserID: %s", userID)
	err := h.urlService.DeleteUrlsByUserID(userID)
	if err != nil {
		logger.GetSugar().Errorf("Error while deleting user URLs: %w", err)
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func getUserIDFromContext(ctx context.Context) string {
	var userID string
	if ctx.Value("userID") != nil {
		userID = ctx.Value("userID").(string)
	} else {
		userID = ""
	}
	return userID
}

func setUserUrlsToJSON(input []*model.Urls) ([]byte, error) {
	output, err := getUrlsJSONFromUrls(input)
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getUrlsJSONFromUrls(urls []*model.Urls) ([]OutputUserUrlsJson, error) {
	output := make([]OutputUserUrlsJson, 0, len(urls))
	for _, u := range urls {
		short, err := service.MakeShortURL(u.GetShortURL())
		if err != nil {
			return output, err
		}
		output = append(output, OutputUserUrlsJson{
			ShortURL:    short,
			OriginalURL: u.GetOriginal(),
		})
	}

	return output, nil
}
