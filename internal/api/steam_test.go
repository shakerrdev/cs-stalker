package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeMockClient(t *testing.T, responseBody string, statusCode int) *http.Client {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		io.WriteString(w, responseBody)
	}))
	t.Cleanup(server.Close)

	return &http.Client{
		Transport: rewriteHostTransport(server.URL),
	}
}

func TestResolveVanityUrl(t *testing.T) {
	mockResp := `{"response":{"steamid":"76561198056395137","success":1}}`
	client := makeMockClient(t, mockResp, http.StatusOK)
	setSteamClient(client)
	t.Setenv("STEAMAPIKEY", "mock")

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
	client := makeMockClient(t, mockResp, http.StatusOK)
	setSteamClient(client)
	t.Setenv("STEAMAPIKEY", "mock")

	_, err := resolveVanityUrl("unknown_user")
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestUserSummary(t *testing.T) {
	mockResp := `{"response":{"players":[{"steamid":"76561198056395137","communityvisibilitystate":3,"personaname":"TestUser","avatarfull":"http://example.com/avatar.jpg"}]}}`
	client := makeMockClient(t, mockResp, http.StatusOK)
	setSteamClient(client)
	t.Setenv("STEAMAPIKEY", "mock")

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
	client := makeMockClient(t, mockResp, http.StatusOK)
	setSteamClient(client)
	t.Setenv("STEAMAPIKEY", "mock")

	_, err := UserSummary("76561198056395137")
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestUserPlaytime(t *testing.T) {
	mockResp := `{"response":{"games":[{"appid":730,"playtime_forever":120}]}}`
	client := makeMockClient(t, mockResp, http.StatusOK)
	setSteamClient(client)
	t.Setenv("STEAMAPIKEY", "mock")

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
	client := makeMockClient(t, mockResp, http.StatusOK)
	setSteamClient(client)
	t.Setenv("STEAMAPIKEY", "mock")

	_, err := UserPlaytime("76561198056395137")
	if err == nil {
		t.Fatal("expected error, got none")
	}
}
