package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/FoPQer/go-shortener/internal/logger"
)

// DBPinger is an interface for database ping operations.
type DBPinger interface {
	Ping(ctx context.Context) error
}

type DBHandler struct {
	db DBPinger
}

func NewDBHandler(db DBPinger) *DBHandler {
	return &DBHandler{db: db}
}

func (h *DBHandler) GetPing(res http.ResponseWriter, req *http.Request) {
	if h.db == nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := h.db.Ping(ctx); err != nil {
		logger.GetSugar().Errorf("Database ping failed: %v", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
