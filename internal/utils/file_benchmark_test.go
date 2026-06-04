package utils

import (
	"path/filepath"
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
)

func BenchmarkWriteToFile(b *testing.B) {
	urls := []*model.Urls{
		model.NewUrls("https://example.com/1", "1"),
		model.NewUrls("https://example.com/2", "2"),
		model.NewUrls("https://example.com/3", "3"),
	}
	filePath := filepath.Join(b.TempDir(), "urls.json")

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := WriteToFile(filePath, urls); err != nil {
			b.Fatalf("WriteToFile failed: %v", err)
		}
	}
}
