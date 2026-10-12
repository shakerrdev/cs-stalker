package api

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func steamResponds(t *testing.T, body string) {
	t.Helper()
	t.Setenv("STEAMAPIKEY", "mock")
	useSteamServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, body)
	})
}

func TestResolveVanityUrl(t *testing.T) {
	mockResp := `{"response":{"steamid":"76561198056395137","success":1}}`
	steamResponds(t, mockResp)

	id, err := resolveVanityUrl("testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "76561198056395137" {
		t.Errorf("expected id 76561198056395137, got %s", id)
	}
}

func TestResolveVanityUrl_Failure(t *testing.T) {
	mockResp := `{"response":{"success":42,"message":"No match"}}`
	steamResponds(t, mockResp)

	_, err := resolveVanityUrl("unknown_user")
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestUserSummary(t *testing.T) {
	mockResp := `{"response":{"players":[{"steamid":"76561198056395137","communityvisibilitystate":3,"personaname":"TestUser","avatarfull":"http://example.com/avatar.jpg"}]}}`
	steamResponds(t, mockResp)

	summary, err := UserSummary("76561198056395137")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Name != "TestUser" {
		t.Errorf("expected name TestUser, got %s", summary.Name)
	}
	if !summary.IsPublic {
		t.Errorf("expected IsPublic to be true")
	}
}

func TestUserSummary_EmptyResponse(t *testing.T) {
	mockResp := `{"response":{"players":[]}}`
	steamResponds(t, mockResp)

	_, err := UserSummary("76561198056395137")
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestUserPlaytime(t *testing.T) {
	t.Setenv("STEAMAPIKEY", "mock")
	useSteamServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "GetOwnedGames") || r.URL.Query().Get("include_played_free_games") != "1" {
			t.Errorf("unexpected request %s", r.URL)
		}
		io.WriteString(w, `{"response":{"games":[{"appid":10,"playtime_forever":30},{"appid":730,"playtime_forever":120}]}}`)
	})

	hours, err := UserPlaytime("76561198056395137")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hours != 2 {
		t.Errorf("expected 2 hours, got %d", hours)
	}
}

func TestUserPlaytime_NotFound(t *testing.T) {
	mockResp := `{"response":{"games":[{"appid":999,"playtime_forever":200}]}}`
	steamResponds(t, mockResp)

	_, err := UserPlaytime("76561198056395137")
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestSteamRequestEscapesParams(t *testing.T) {
	t.Setenv("STEAMAPIKEY", "secret-key")
	var query url.Values
	useSteamServer(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		io.WriteString(w, `{"response":{"steamid":"76561198056395137","success":1}}`)
	})

	if _, err := resolveVanityUrl("name&key=other"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := query.Get("vanityurl"); got != "name&key=other" {
		t.Errorf("vanityurl = %q", got)
	}
	if got := query["key"]; len(got) != 1 || got[0] != "secret-key" {
		t.Errorf("key = %q", got)
	}
}

func TestSteamRequestKeepsKeyOutOfLogsAndErrors(t *testing.T) {
	t.Setenv("STEAMAPIKEY", "secret-key")
	previous := steamClient
	steamClient = &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})}
	t.Cleanup(func() { steamClient = previous })

	var logs bytes.Buffer
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	_, err := UserSummary("76561198056395137")
	if err == nil {
		t.Fatal("expected error, got none")
	}
	if strings.Contains(err.Error(), "secret-key") {
		t.Errorf("error leaks api key: %v", err)
	}
	if strings.Contains(logs.String(), "secret-key") {
		t.Errorf("log leaks api key: %s", logs.String())
	}
}

func TestResolveUrl(t *testing.T) {
	steamResponds(t, `{"response":{"steamid":"76561197960287930","success":1}}`)

	tests := []struct {
		input string
		want  string
	}{
		{"76561198056395137", "76561198056395137"},
		{"  76561198056395137  ", "76561198056395137"},
		{"https://steamcommunity.com/profiles/76561198056395137", "76561198056395137"},
		{"https://steamcommunity.com/profiles/76561198056395137/", "76561198056395137"},
		{"https://steamcommunity.com/profiles/76561198056395137/?l=german", "76561198056395137"},
		{"steamcommunity.com/profiles/76561198056395137", "76561198056395137"},
		{"https://steamcommunity.com/id/gabelogannewell", "76561197960287930"},
		{"https://steamcommunity.com/id/gabelogannewell/#top", "76561197960287930"},
		{"gabelogannewell", "76561197960287930"},
	}
	for _, tt := range tests {
		got, err := ResolveUrl(tt.input)
		if err != nil {
			t.Errorf("ResolveUrl(%q): unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ResolveUrl(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestResolveUrl_Invalid(t *testing.T) {
	for _, input := range []string{
		"",
		"https://steamcommunity.com/profiles/not-an-id",
		"https://steamcommunity.com/profiles/123",
		"https://steamcommunity.com/groups/somegroup",
		"https://example.com/",
	} {
		if got, err := ResolveUrl(input); err == nil {
			t.Errorf("ResolveUrl(%q) = %q, want an error", input, got)
		}
	}
}
