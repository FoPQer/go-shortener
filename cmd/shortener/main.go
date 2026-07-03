package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/FoPQer/go-shortener/internal/auth"
	"github.com/FoPQer/go-shortener/internal/config"
	"github.com/FoPQer/go-shortener/internal/config/db"
	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/events"
	"github.com/FoPQer/go-shortener/internal/handlers"
	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/middlewares"
	repoFactory "github.com/FoPQer/go-shortener/internal/repository/factory"
	"github.com/FoPQer/go-shortener/internal/routes"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/FoPQer/go-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	initMsg()

	flags.ParseFlags()
	logger.InitLogger()

	// Load configuration from file if specified
	configFilePath := os.Getenv("CONFIG")
	if configFilePath == "" {
		configFilePath = flags.GetFlagConfigFile()
	}
	if configFilePath != "" {
		if _, err := config.LoadConfig(configFilePath); err != nil {
			logger.GetSugar().Warnf("Failed to load configuration file: %v", err)
		}
	}

	pgxConf, err := db.InitPgsql()
	if errors.Is(err, db.ErrConnNotFound) {
		logger.GetSugar().Infoln("Database connection string not found, using file or memory repository")
	} else if err != nil {
		logger.GetSugar().Fatalf("Error initializing database: %s", err.Error())
	}
	if pgxConf.GetDBConn() != nil {
		defer pgxConf.GetDBConn().Close()
	}

	factory := repoFactory.NewRepositoryFactory(pgxConf.GetDBConn(), service.GetFileStoragePath())

	urlRepo, err := factory.CreateUrlsRepository()
	if err != nil {
		logger.GetSugar().Fatal(err)
	}

	userRepo, err := factory.CreateUserRepository()
	if err != nil {
		logger.GetSugar().Fatal(err)
	}

	urlService := service.NewURLService(urlRepo)
	jsonService := service.NewJSONService()
	userService := service.NewUserService(userRepo)
	statService := service.NewStatService(urlService, userService)
	claimsService := auth.NewClaimsService()

	authMiddleware := middlewares.NewAuthMiddleware(userService, claimsService)
	trustedMiddleware := middlewares.NewTrustedMiddleware(service.GetTrustedSubnet())

	auditFilePath := service.GetAuditFile()
	auditURLPath := service.GetAuditURL()
	var auditPublisher events.Publisher
	if auditFilePath == "" && auditURLPath == "" {
		logger.GetSugar().Infoln("No audit destination specified, skipping audit setup")
	} else {
		auditBus := events.NewAuditBus(100)
		auditPublisher = auditBus

		if auditFilePath != "" {
			auditFile := events.NewAuditFile(1, auditFilePath)
			auditBus.AddSubscriber(auditFile)
			logger.GetSugar().Infoln("Audit file successfully setup")
		}

		if auditURLPath != "" {
			auditURL := events.NewAuditURL(1, auditURLPath)
			auditBus.AddSubscriber(auditURL)
			logger.GetSugar().Infoln("Audit url successfully setup")
		}
	}

	handler := handlers.NewHandler(urlService, jsonService, userService, statService, auditPublisher)
	dbHandler := handlers.NewDBHandler(pgxConf.GetDBConn())

	r := chi.NewRouter()
	routes.InitWebRoutes(r, handler, dbHandler, authMiddleware, trustedMiddleware)
	r.Mount("/debug", chiMiddleware.Profiler())

	srv := &http.Server{
		Addr:    service.GetRunAddr(),
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	serverErr := make(chan error, 1)

	if service.GetHTTPs() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.GetSugar().Fatal(err)
		}

		certPath := filepath.Join(homeDir, "cert.pem")
		privateKeyPath := filepath.Join(homeDir, "private.pem")

		certExists := true
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			certExists = false
		}
		keyExists := true
		if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
			keyExists = false
		}

		if !certExists || !keyExists {
			logger.GetSugar().Infoln("Certificate files not found, generating self-signed certificates")
			if err := utils.GenerateCertificate(); err != nil {
				logger.GetSugar().Fatal(err)
			}
		}

		logger.GetSugar().Infoln("Starting server with HTTPS")
		go func() {
			serverErr <- srv.ListenAndServeTLS(certPath, privateKeyPath)
		}()
	} else {
		logger.GetSugar().Infoln("Starting server with HTTP")
		go func() {
			serverErr <- srv.ListenAndServe()
		}()
	}

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.GetSugar().Fatalf("Server error: %s", err)
		}
	case <-ctx.Done():
		logger.GetSugar().Infoln("Received shutdown signal, shutting down gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.GetSugar().Errorf("Graceful shutdown failed: %s", err)
		} else {
			logger.GetSugar().Infoln("Server shut down successfully")
		}

		if auditBus, ok := auditPublisher.(*events.AuditBus); ok {
			auditBus.Close()
			logger.GetSugar().Infoln("Audit bus closed")
		}
	}
}

func initMsg() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
