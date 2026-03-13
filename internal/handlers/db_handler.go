package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/FoPQer/go-shortener/internal/config/db"
	"github.com/FoPQer/go-shortener/internal/logger"
)

type DBHandler struct {
	db *db.PgxConf
}

func NewDBHandler(db *db.PgxConf) *DBHandler {
	return &DBHandler{db: db}
}

func (h *DBHandler) GetPing(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := h.db.GetDBConn().Ping(ctx); err != nil {
		logger.GetSugar().Errorf("Database ping failed: %v", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}