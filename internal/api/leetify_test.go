package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// sampleProfile is a stubbed response for testing
var sampleProfile = Profile{
	Meta: struct {
		Name           string   `json:"name"`
		Steam64ID      string   `json:"steam64Id"`
		SteamAvatarURL string   `json:"steamAvatarUrl"`
		FaceitNickname string   `json:"faceitNickname"`
		PlatformBans   []string `json:"platformBans"`
	}{
		Name:           "TestUser",
		Steam64ID:      "123456789",
		SteamAvatarURL: "http://example.com/avatar.jpg",
		FaceitNickname: "TestFaceit",
		PlatformBans:   []string{},
	},
}

func TestLeetifyProfile_Success(t *testing.T) {
	// Create a fake HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return success only if URL ends with "/id/testid"
		if r.URL.Path == "/api/profile/id/testid" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(sampleProfile)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	// Patch the base URLs used in LeetifyProfile
	originalEnv := os.Getenv("LEETIFY_AUTH_TOKEN")
	defer os.Setenv("LEETIFY_AUTH_TOKEN", originalEnv)
	os.Setenv("LEETIFY_AUTH_TOKEN", "mocked-token")

	// Set up the custom HTTP client and set it for testing
	leetifyClient := &http.Client{
		Transport: rewriteHostTransport(ts.URL),
	}
	setLeetifyClient(leetifyClient) // Use the setter function to apply the custom client

	// Call the LeetifyProfile function
	profile, err := LeetifyProfile("testid")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if profile.Meta.Name != "TestUser" {
		t.Errorf("Expected profile.Meta.Name = TestUser, got %s", profile.Meta.Name)
	}
}

func TestLeetifyProfile_InvalidToken(t *testing.T) {
	// Create a fake HTTP server that simulates the API rejecting the token
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized) // Unauthorized error for invalid token
	}))
	defer ts.Close()

	// Set invalid token environment variable
	originalEnv := os.Getenv("LEETIFY_AUTH_TOKEN")
	defer os.Setenv("LEETIFY_AUTH_TOKEN", originalEnv)
	os.Setenv("LEETIFY_AUTH_TOKEN", "invalid-token")

	// Set up the custom HTTP client and set it for testing
	leetifyClient := &http.Client{
		Transport: rewriteHostTransport(ts.URL),
	}
	setLeetifyClient(leetifyClient)

	// Call the LeetifyProfile function
	_, err := LeetifyProfile("testid")
	if err == nil {
		t.Fatalf("Expected error due to invalid token, but got none")
	}
}

func TestLeetifyProfile_NotFound(t *testing.T) {
	// Create a fake HTTP server that simulates a "not found" response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}))
	defer ts.Close()

	// Set the token environment variable
	originalEnv := os.Getenv("LEETIFY_AUTH_TOKEN")
	defer os.Setenv("LEETIFY_AUTH_TOKEN", originalEnv)
	os.Setenv("LEETIFY_AUTH_TOKEN", "mocked-token")

	// Set up the custom HTTP client and set it for testing
	leetifyClient := &http.Client{
		Transport: rewriteHostTransport(ts.URL),
	}
	setLeetifyClient(leetifyClient)

	// Call the LeetifyProfile function
	_, err := LeetifyProfile("unknownid")
	if err == nil {
		t.Fatalf("Expected error due to profile not found, but got none")
	}
}

func TestLeetifyProfile_APIError(t *testing.T) {
	// Create a fake HTTP server that simulates an API error (500 Internal Server Error)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
	}))
	defer ts.Close()

	// Set the token environment variable
	originalEnv := os.Getenv("LEETIFY_AUTH_TOKEN")
	defer os.Setenv("LEETIFY_AUTH_TOKEN", originalEnv)
	os.Setenv("LEETIFY_AUTH_TOKEN", "mocked-token")

	// Set up the custom HTTP client and set it for testing
	leetifyClient := &http.Client{
		Transport: rewriteHostTransport(ts.URL),
	}
	setLeetifyClient(leetifyClient)

	// Call the LeetifyProfile function
	_, err := LeetifyProfile("testid")
	if err == nil {
		t.Fatalf("Expected error due to API error, but got none")
	}
}
