package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/user"
)

type mockUserRepository struct {
	findByIDFn func(ctx context.Context, id string) (*model.User, error)
	saveFn     func(ctx context.Context, user *model.User) (string, error)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}

	return model.NewUser(id), nil
}

func (m *mockUserRepository) Save(ctx context.Context, userModel *model.User) (string, error) {
	if m.saveFn != nil {
		return m.saveFn(ctx, userModel)
	}

	return userModel.GetID(), nil
}

var _ user.UserRepository = (*mockUserRepository)(nil)

func benchmarkStringGetter(b *testing.B, envKey, envValue string, getter func() string) {
	b.Helper()
	b.Setenv(envKey, envValue)
	resetConfigCache()
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		value := getter()
		if value == "" {
			b.Fatal("expected non-empty result")
		}
	}
}

func BenchmarkGetRunAddr(b *testing.B) {
	benchmarkStringGetter(b, "SERVER_ADDRESS", "127.0.0.1:8080", GetRunAddr)
}

func BenchmarkGetBasePrefix(b *testing.B) {
	benchmarkStringGetter(b, "BASE_URL", "short", GetBasePrefix)
}

func BenchmarkGetFileStoragePath(b *testing.B) {
	benchmarkStringGetter(b, "FILE_STORAGE_PATH", `C:\bench\storage.json`, GetFileStoragePath)
}

func BenchmarkGetDatabaseDSN(b *testing.B) {
	benchmarkStringGetter(b, "DATABASE_DSN", "postgres://bench", GetDatabaseDSN)
}

func BenchmarkGetSecretKey(b *testing.B) {
	benchmarkStringGetter(b, "SECRET_KEY", "bench-secret", GetSecretKey)
}

func BenchmarkGetAuditFile(b *testing.B) {
	benchmarkStringGetter(b, "AUDIT_FILE", `  "C:\bench\audit file.json"  `, GetAuditFile)
}

func BenchmarkGetAuditURL(b *testing.B) {
	benchmarkStringGetter(b, "AUDIT_URL", "https://example.com/audit", GetAuditURL)
}

func BenchmarkJSONService_GetURLFromJSON(b *testing.B) {
	svc := NewJSONService()
	input := []byte(`{"url":"https://example.com"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := svc.GetURLFromJSON(input)
		if err != nil {
			b.Fatalf("GetURLFromJSON failed: %v", err)
		}
		if got == "" {
			b.Fatal("expected non-empty URL")
		}
	}
}

func BenchmarkJSONService_SetURLToJSON(b *testing.B) {
	svc := NewJSONService()

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := svc.SetURLToJSON("https://example.com")
		if err != nil {
			b.Fatalf("SetURLToJSON failed: %v", err)
		}
		if len(got) == 0 {
			b.Fatal("expected non-empty JSON output")
		}
	}
}

func BenchmarkJSONService_GetUrlsFromInputBatchJSON(b *testing.B) {
	svc := NewJSONService()
	input := make([]InputJSONBatchURL, 128)
	for i := range input {
		input[i] = InputJSONBatchURL{
			CorrelationID: "id-" + string(rune('a'+(i%26))),
			OriginalURL:   "https://example.com/item",
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := svc.GetUrlsFromInputBatchJSON(input)
		if err != nil {
			b.Fatalf("GetUrlsFromInputBatchJSON failed: %v", err)
		}
		if len(got) != len(input) {
			b.Fatalf("expected %d URLs, got %d", len(input), len(got))
		}
	}
}

func BenchmarkJSONService_GetBatchURLFromJSON(b *testing.B) {
	svc := NewJSONService()
	inputStruct := make([]InputJSONBatchURL, 128)
	for i := range inputStruct {
		inputStruct[i] = InputJSONBatchURL{
			CorrelationID: "id-" + string(rune('a'+(i%26))),
			OriginalURL:   "https://example.com/item",
		}
	}
	input, err := json.Marshal(inputStruct)
	if err != nil {
		b.Fatalf("failed to prepare benchmark input: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := svc.GetBatchURLFromJSON(input)
		if err != nil {
			b.Fatalf("GetBatchURLFromJSON failed: %v", err)
		}
		if len(got) != len(inputStruct) {
			b.Fatalf("expected %d URLs, got %d", len(inputStruct), len(got))
		}
	}
}

func BenchmarkJSONService_GetOutputBatchJSONFromUrls(b *testing.B) {
	b.Setenv("SERVER_ADDRESS", "127.0.0.1:8080")
	b.Setenv("BASE_URL", "short")
	resetConfigCache()

	svc := NewJSONService()
	input := make([]*model.Urls, 128)
	for i := range input {
		input[i] = model.NewUrls("https://example.com/item", "id-"+string(rune('a'+(i%26))))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := svc.GetOutputBatchJSONFromUrls(input)
		if err != nil {
			b.Fatalf("GetOutputBatchJSONFromUrls failed: %v", err)
		}
		if len(got) != len(input) {
			b.Fatalf("expected %d output items, got %d", len(input), len(got))
		}
	}
}

func BenchmarkJSONService_SetBatchURLToJSON(b *testing.B) {
	b.Setenv("SERVER_ADDRESS", "127.0.0.1:8080")
	b.Setenv("BASE_URL", "short")
	resetConfigCache()

	svc := NewJSONService()
	input := make([]*model.Urls, 128)
	for i := range input {
		input[i] = model.NewUrls("https://example.com/item", "id-"+string(rune('a'+(i%26))))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := svc.SetBatchURLToJSON(input)
		if err != nil {
			b.Fatalf("SetBatchURLToJSON failed: %v", err)
		}
		if len(got) == 0 {
			b.Fatal("expected non-empty JSON output")
		}
	}
}

func BenchmarkUserService_Get(b *testing.B) {
	repo := &mockUserRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return model.NewUser(id), nil
		},
	}
	svc := NewUserService(repo)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := svc.Get(context.Background(), "user-1")
		if err != nil {
			b.Fatalf("Get failed: %v", err)
		}
		if got == nil {
			b.Fatal("expected user result")
		}
	}
}

func BenchmarkUserService_Create(b *testing.B) {
	repo := &mockUserRepository{
		saveFn: func(ctx context.Context, userModel *model.User) (string, error) {
			return userModel.GetID(), nil
		},
	}
	svc := NewUserService(repo)
	input := model.NewUser("user-1")

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := svc.Create(context.Background(), input)
		if err != nil {
			b.Fatalf("Create failed: %v", err)
		}
		if got == "" {
			b.Fatal("expected user ID")
		}
	}
}
