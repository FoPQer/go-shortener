package factory

import (
	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
	urlsDB "github.com/FoPQer/go-shortener/internal/repository/urls/db"
	urlsFile "github.com/FoPQer/go-shortener/internal/repository/urls/file"
	urlsMemory "github.com/FoPQer/go-shortener/internal/repository/urls/memory"

	"github.com/FoPQer/go-shortener/internal/repository/user"
	userDB "github.com/FoPQer/go-shortener/internal/repository/user/db"

	userMemory "github.com/FoPQer/go-shortener/internal/repository/user/memory"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RepositoryFactory builds repository implementations based on runtime configuration.
type RepositoryFactory struct {
	conn     *pgxpool.Pool
	filePath string
}

// NewRepositoryFactory constructs a RepositoryFactory with database and file settings.
func NewRepositoryFactory(conn *pgxpool.Pool, filePath string) *RepositoryFactory {
	return &RepositoryFactory{conn: conn, filePath: filePath}
}

// CreateUrlsRepository returns a URL repository implementation.
//
// Priority: database repository when DB connection is available,
// file repository when file path is configured, otherwise in-memory repository.
func (f *RepositoryFactory) CreateUrlsRepository() (urls.Repository, error) {
	var repo urls.Repository
	if f.conn != nil {
		repo = urlsDB.NewRepository(f.conn)
		logger.GetSugar().Infoln("Working with database repository")
	} else if f.filePath != "" {
		repo = urlsFile.NewRepository(f.filePath)
		logger.GetSugar().Infoln("Working with file repository")
	} else {
		repo = urlsMemory.NewRepository()
		logger.GetSugar().Infoln("Working with memory repository")
	}

	return repo, nil
}

// CreateUserRepository returns a user repository implementation.
//
// Priority: database repository when DB connection is available,
// otherwise in-memory repository.
func (f *RepositoryFactory) CreateUserRepository() (user.UserRepository, error) {
	var repo user.UserRepository
	if f.conn != nil {
		repo = userDB.NewRepository(f.conn)
		logger.GetSugar().Infoln("Working with database repository")
	} else {
		repo = userMemory.NewRepository()
		logger.GetSugar().Infoln("Working with memory repository")
	}

	return repo, nil
}
