package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequirePassword(t *testing.T) {
	handler := RequirePassword("hunter2", newTestRouter(t))

	tests := []struct {
		name     string
		target   string
		user     string
		password string
		auth     bool
		status   int
	}{
		{"no credentials", "/", "", "", false, http.StatusUnauthorized},
		{"wrong password", "/", "me", "hunter3", true, http.StatusUnauthorized},
		{"right password", "/", "me", "hunter2", true, http.StatusOK},
		{"any user name", "/", "", "hunter2", true, http.StatusOK},
		{"api needs it too", "/api/resolveurl?url=76561198056395137", "", "", false, http.StatusUnauthorized},
		{"health check stays open", "/healthz", "", "", false, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.target, nil)
			if tt.auth {
				req.SetBasicAuth(tt.user, tt.password)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}
			if tt.status == http.StatusUnauthorized && rec.Header().Get("WWW-Authenticate") == "" {
				t.Error("missing WWW-Authenticate header, browsers won't ask for the password")
			}
		})
	}
}
