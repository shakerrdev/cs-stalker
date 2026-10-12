package routes

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestUnknownPathIsNotFound(t *testing.T) {
	if rec := get(newTestRouter(t), "/does-not-exist"); rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestProfileLookupShowsErrors(t *testing.T) {
	form := url.Values{"url": {"https://example.com/not-a-profile"}}
	req := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	newTestRouter(t).ServeHTTP(rec, req)

	if !strings.Contains(rec.Body.String(), "No Steam profile found") {
		t.Errorf("body = %q", rec.Body.String())
	}
}

func TestHealthz(t *testing.T) {
	rec := get(newTestRouter(t), "/healthz")
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Errorf("GET /healthz = %d %q", rec.Code, rec.Body.String())
	}
}
