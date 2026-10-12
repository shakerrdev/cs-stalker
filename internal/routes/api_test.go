package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiResolveUrlHandler(t *testing.T) {
	tests := []struct {
		name   string
		target string
		status int
		key    string
		value  string
	}{
		{"steam64 id", "/api/resolveurl?url=76561198056395137", http.StatusOK, "steam64id", "76561198056395137"},
		{"missing url", "/api/resolveurl", http.StatusBadRequest, "error", "missing url parameter"},
		{"not a profile", "/api/resolveurl?url=https://example.com/", http.StatusNotFound, "error", "could not resolve a steam profile"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			apiResolveUrlHandler(rec, httptest.NewRequest(http.MethodGet, tt.target, nil))

			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("content type = %q", ct)
			}
			var body map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body[tt.key] != tt.value {
				t.Errorf("body = %v, want %s=%q", body, tt.key, tt.value)
			}
		})
	}
}
