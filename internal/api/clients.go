package api

import (
	"fmt"
	"net/http"
	"os"
)

// authTransport adds an Authorization header to the request.
type authTransport struct {
	Transport http.RoundTripper
}

func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	authToken := os.Getenv("LEETIFY_AUTH_TOKEN")
	if authToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}
	return a.Transport.RoundTrip(req)
}

// roundTripperFunc is a helper for wrapping RoundTrip logic in tests.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// rewriteHostTransport rewrites request destinations to the test server host.
func rewriteHostTransport(testServerURL string) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = "http"
		req.URL.Host = testServerURL[len("http://"):]
		return http.DefaultTransport.RoundTrip(req)
	})
}

// Clients used by external APIs. Can be swapped for testing.
var (
	steamClient   *http.Client
	leetifyClient *http.Client
)

func setSteamClient(c *http.Client) {
	if c != nil {
		steamClient = c
	}
}

func setLeetifyClient(c *http.Client) {
	if c != nil {
		leetifyClient = c
	}
}

func init() {
	leetifyClient = &http.Client{
		Transport: &authTransport{
			Transport: http.DefaultTransport,
		},
	}

	steamClient = http.DefaultClient
}
