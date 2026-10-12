package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	router, err := NewRouter("../..")
	if err != nil {
		t.Fatal(err)
	}
	return router
}

func get(router http.Handler, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestPages(t *testing.T) {
	router := newTestRouter(t)
	tests := []struct {
		target   string
		contains string
	}{
		{"/", `hx-post="/profile"`},
		{"/minion", `id="minion"`},
		{"/static/index.css", "box-sizing"},
	}
	for _, tt := range tests {
		rec := get(router, tt.target)
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s: status = %d", tt.target, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), tt.contains) {
			t.Errorf("GET %s: body does not contain %q", tt.target, tt.contains)
		}
	}
}
