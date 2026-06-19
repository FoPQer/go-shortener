package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/FoPQer/go-shortener/internal/auth"
	"github.com/FoPQer/go-shortener/internal/handlers"
	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/middlewares"
	urlsFile "github.com/FoPQer/go-shortener/internal/repository/urls/file"
	urlMemory "github.com/FoPQer/go-shortener/internal/repository/urls/memory"
	userMemory "github.com/FoPQer/go-shortener/internal/repository/user/memory"
	"github.com/FoPQer/go-shortener/internal/routes"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

var initExampleLoggerOnce sync.Once

func initExampleLogger() {
	initExampleLoggerOnce.Do(func() {
		if err := logger.InitLogger(); err != nil {
			panic(err)
		}
	})
}

func newExampleRouter() *chi.Mux {
	urlService := service.NewURLService(urlMemory.NewRepository())
	jsonService := service.NewJSONService()
	userService := service.NewUserService(userMemory.NewRepository())
	statService := service.NewStatService(urlService, userService)
	claimsService := auth.NewClaimsService()

	handler := handlers.NewHandler(urlService, jsonService, userService, statService, nil)
	dbHandler := handlers.NewDBHandler(nil)
	authMiddleware := middlewares.NewAuthMiddleware(userService, claimsService)
	trustedMiddleware := middlewares.NewTrustedMiddleware("")

	r := chi.NewRouter()
	routes.InitWebRoutes(r, handler, dbHandler, authMiddleware, trustedMiddleware)
	return r
}

func newExampleRouterWithFileURLRepo(filePath string) *chi.Mux {
	urlService := service.NewURLService(urlsFile.NewRepository(filePath))
	jsonService := service.NewJSONService()
	userService := service.NewUserService(userMemory.NewRepository())
	statService := service.NewStatService(urlService, userService)
	claimsService := auth.NewClaimsService()

	handler := handlers.NewHandler(urlService, jsonService, userService, statService, nil)
	dbHandler := handlers.NewDBHandler(nil)
	authMiddleware := middlewares.NewAuthMiddleware(userService, claimsService)
	trustedMiddleware := middlewares.NewTrustedMiddleware("")

	r := chi.NewRouter()
	routes.InitWebRoutes(r, handler, dbHandler, authMiddleware, trustedMiddleware)
	return r
}

func extractShortID(shortURL string) string {
	u, err := url.Parse(shortURL)
	if err != nil {
		panic(err)
	}

	id := path.Base(u.Path)
	if id == "." || id == "/" || id == "" {
		id = u.Host
	}

	if id == "" {
		panic("empty short id")
	}

	return id
}

func Example_plainTextShortenAndFollow() {
	initExampleLogger()
	r := newExampleRouter()

	postReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com/path"))
	postRec := httptest.NewRecorder()
	r.ServeHTTP(postRec, postReq)
	cookies := postRec.Result().Cookies()
	if len(cookies) == 0 {
		panic("auth cookie was not set")
	}

	shortURL := strings.TrimSpace(postRec.Body.String())
	id := extractShortID(shortURL)

	getReq := httptest.NewRequest(http.MethodGet, "/"+id, nil)
	getReq.AddCookie(cookies[0])
	getRec := httptest.NewRecorder()
	r.ServeHTTP(getRec, getReq)

	fmt.Println(postRec.Code)
	fmt.Println(getRec.Code)
	fmt.Println(getRec.Header().Get("Location"))
	// Output:
	// 201
	// 307
	// https://example.com/path
}

func Example_jsonShorten() {
	initExampleLogger()
	r := newExampleRouter()

	body := `{"url":"https://example.com/docs"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	responseBody, err := io.ReadAll(rec.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(rec.Code)
	fmt.Println(strings.Contains(string(responseBody), "result"))
	// Output:
	// 201
	// true
}

func Example_userURLsList() {
	initExampleLogger()

	tmp, err := os.CreateTemp("", "go-shortener-example-*.json")
	if err != nil {
		panic(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	r := newExampleRouterWithFileURLRepo(tmp.Name())

	createBody := `{"url":"https://example.com/user-page"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	r.ServeHTTP(createRec, createReq)

	cookies := createRec.Result().Cookies()
	if len(cookies) == 0 {
		panic("auth cookie was not set")
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	listReq.AddCookie(cookies[0])
	listRec := httptest.NewRecorder()
	r.ServeHTTP(listRec, listReq)

	listBody, err := io.ReadAll(listRec.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(listRec.Code)
	fmt.Println(strings.Contains(string(listBody), "original_url"))
	// Output:
	// 200
	// true
}

func Example_deleteUserURLs() {
	initExampleLogger()
	r := newExampleRouter()

	createBody := `{"url":"https://example.com/to-delete"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	r.ServeHTTP(createRec, createReq)

	cookies := createRec.Result().Cookies()
	if len(cookies) == 0 {
		panic("auth cookie was not set")
	}

	var shortenResp struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(createRec.Body).Decode(&shortenResp); err != nil {
		panic(err)
	}
	shortID := extractShortID(shortenResp.Result)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBufferString(fmt.Sprintf("[\"%s\"]", shortID)))
	deleteReq.Header.Set("Content-Type", "application/json")
	deleteReq.AddCookie(cookies[0])
	deleteRec := httptest.NewRecorder()
	r.ServeHTTP(deleteRec, deleteReq)

	followReq := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
	followReq.AddCookie(cookies[0])
	followRec := httptest.NewRecorder()
	r.ServeHTTP(followRec, followReq)

	fmt.Println(deleteRec.Code)
	fmt.Println(followRec.Code)
	// Output:
	// 202
	// 410
}
