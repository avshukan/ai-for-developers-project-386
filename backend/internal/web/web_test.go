package web_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/avshukan/ai-for-developers-project-386/backend/internal/web"
)

func TestSPAHandler(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "index.html"), "<!doctype html><title>app</title>")
	writeFile(t, filepath.Join(dir, "assets", "app.js"), "console.log('hi')")

	h := web.SPAHandler(dir)

	tests := []struct {
		name     string
		path     string
		wantBody string
	}{
		{"root serves index", "/", "<!doctype html><title>app</title>"},
		{"existing asset is served", "/assets/app.js", "console.log('hi')"},
		{"client route falls back to index", "/bookings/123", "<!doctype html><title>app</title>"},
		{"missing asset falls back to index", "/assets/missing.js", "<!doctype html><title>app</title>"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}
			if got := rec.Body.String(); got != tc.wantBody {
				t.Fatalf("body = %q, want %q", got, tc.wantBody)
			}
		})
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
