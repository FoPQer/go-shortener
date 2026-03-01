package factory

import (
	"context"
	"log"

	"github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/FoPQer/go-shortener/internal/repository/urls/db"
	"github.com/FoPQer/go-shortener/internal/repository/urls/file"
	"github.com/FoPQer/go-shortener/internal/repository/urls/memory"
	"github.com/jackc/pgx/v5"
)

type RepositoryFactory struct {
	conn *pgx.Conn
	filePath string
}

func NewRepositoryFactory(conn *pgx.Conn, filePath string) *RepositoryFactory {
	return &RepositoryFactory{conn: conn, filePath: filePath}
}

func (f *RepositoryFactory) CreateUrlsRepository(ctx context.Context) (urls.Repository, error) {
	var repo urls.Repository
	if f.conn != nil {
		repo = db.NewRepository(f.conn)
		log.Println("Working with database repository")
	} else if f.filePath != "" {
		repo = file.NewRepository(f.filePath)
		log.Println("Working with file repository")
	} else {
		repo = memory.NewRepository()
		log.Println("Working with memory repository")
	}

	return repo, nil
}