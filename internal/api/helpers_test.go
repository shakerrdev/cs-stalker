package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func rewriteHostTransport(serverURL string) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = "http"
		req.URL.Host = strings.TrimPrefix(serverURL, "http://")
		return http.DefaultTransport.RoundTrip(req)
	})
}

func useSteamServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	previous := steamClient
	steamClient = &http.Client{Transport: rewriteHostTransport(server.URL)}
	t.Cleanup(func() { steamClient = previous })
}

func useLeetifyServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	previous := leetifyClient
	leetifyClient = &http.Client{Transport: rewriteHostTransport(server.URL)}
	t.Cleanup(func() { leetifyClient = previous })
}
