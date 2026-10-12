package api

import (
	"fmt"
	"net/http"
	"os"
)

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

var (
	steamClient   = http.DefaultClient
	leetifyClient = &http.Client{Transport: &authTransport{Transport: http.DefaultTransport}}
)
