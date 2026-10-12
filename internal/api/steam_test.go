package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
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
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
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
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestPlayerSummaries(t *testing.T) {
	t.Setenv("STEAMAPIKEY", "mock")
	var asked []string
	useSteamServer(t, func(w http.ResponseWriter, r *http.Request) {
		asked = append(asked, r.URL.Query().Get("steamids"))
		io.WriteString(w, `{"response":{"players":[
			{"steamid":"76561198000000001","personaname":"one","profileurl":"https://steamcommunity.com/id/one/","avatarmedium":"m.jpg","avatarfull":"f.jpg",
			 "communityvisibilitystate":3,"personastate":1,"gameextrainfo":"Counter-Strike 2","timecreated":1331000000,"loccountrycode":"DE"},
			{"steamid":"76561198000000002","personaname":"two","communityvisibilitystate":1}
		]}}`)
	})

	summaries, err := PlayerSummaries([]string{"76561198000000001", "76561198000000002"})
	if err != nil {
		t.Fatal(err)
	}
	if len(asked) != 1 || asked[0] != "76561198000000001,76561198000000002" {
		t.Errorf("requests = %q, want one batched request", asked)
	}
	one := summaries["76561198000000001"]
	want := Summary{
		Steam64id: "76561198000000001", Name: "one", Avatar: "f.jpg", AvatarMedium: "m.jpg",
		ProfileURL: "https://steamcommunity.com/id/one/", IsPublic: true, PersonaState: 1,
		Game: "Counter-Strike 2", Created: time.Unix(1331000000, 0).UTC(), CountryCode: "DE",
	}
	if one != want {
		t.Errorf("summary = %+v\nwant %+v", one, want)
	}
	if two := summaries["76561198000000002"]; two.IsPublic || !two.Created.IsZero() {
		t.Errorf("private summary = %+v", two)
	}
}

func TestPlayerSummaries_BatchesOf100(t *testing.T) {
	t.Setenv("STEAMAPIKEY", "mock")
	var sizes []int
	useSteamServer(t, func(w http.ResponseWriter, r *http.Request) {
		sizes = append(sizes, len(strings.Split(r.URL.Query().Get("steamids"), ",")))
		io.WriteString(w, `{"response":{"players":[]}}`)
	})

	ids := make([]string, 150)
	for i := range ids {
		ids[i] = fmt.Sprintf("765611980%08d", i)
	}
	if _, err := PlayerSummaries(ids); err != nil {
		t.Fatal(err)
	}
	if len(sizes) != 2 || sizes[0] != 100 || sizes[1] != 50 {
		t.Errorf("batch sizes = %v, want [100 50]", sizes)
	}
}

func TestCSPlaytime(t *testing.T) {
	t.Setenv("STEAMAPIKEY", "mock")
	useSteamServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "GetOwnedGames") || r.URL.Query().Get("include_played_free_games") != "1" {
			t.Errorf("unexpected request %s", r.URL)
		}
		io.WriteString(w, `{"response":{"games":[{"appid":10,"playtime_forever":30},{"appid":730,"playtime_forever":12345,"playtime_2weeks":600,"rtime_last_played":1760000000}]}}`)
	})

	playtime, err := CSPlaytime("76561198056395137")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := Playtime{TotalMinutes: 12345, RecentMinutes: 600, LastPlayed: time.Unix(1760000000, 0).UTC()}
	if playtime != want {
		t.Errorf("playtime = %+v, want %+v", playtime, want)
	}
}

func TestCSPlaytime_NeverPlayed(t *testing.T) {
	steamResponds(t, `{"response":{"game_count":1,"games":[{"appid":999,"playtime_forever":200}]}}`)
	if _, err := CSPlaytime("76561198056395137"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestCSPlaytime_Hidden(t *testing.T) {
	steamResponds(t, `{"response":{}}`)
	if _, err := CSPlaytime("76561198056395137"); !errors.Is(err, ErrHidden) {
		t.Fatalf("err = %v, want ErrHidden", err)
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
		if got, err := ResolveUrl(input); !errors.Is(err, ErrNotFound) {
			t.Errorf("ResolveUrl(%q) = %q, %v; want ErrNotFound", input, got, err)
		}
	}
}

func TestPlayerBans(t *testing.T) {
	tests := []struct {
		name string
		body string
		want Bans
		any  bool
	}{
		{"clean", `{"players":[{"SteamId":"76561198056395137","CommunityBanned":false,"VACBanned":false,"NumberOfVACBans":0,"DaysSinceLastBan":0,"NumberOfGameBans":0,"EconomyBan":"none"}]}`, Bans{}, false},
		{"banned", `{"players":[{"SteamId":"76561198056395137","CommunityBanned":true,"VACBanned":true,"NumberOfVACBans":2,"DaysSinceLastBan":340,"NumberOfGameBans":1,"EconomyBan":"probation"}]}`,
			Bans{VACBans: 2, GameBans: 1, DaysSinceLastBan: 340, CommunityBanned: true, TradeBan: "probation"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			steamResponds(t, tt.body)
			bans, err := PlayerBans("76561198056395137")
			if err != nil {
				t.Fatal(err)
			}
			if bans != tt.want || bans.Any() != tt.any {
				t.Errorf("bans = %+v (any %v), want %+v (any %v)", bans, bans.Any(), tt.want, tt.any)
			}
		})
	}
}

func TestSteamLevel(t *testing.T) {
	steamResponds(t, `{"response":{"player_level":42}}`)
	if level, err := SteamLevel("76561198056395137"); err != nil || level != 42 {
		t.Errorf("SteamLevel = %d, %v; want 42", level, err)
	}
}

func TestSteamLevel_Hidden(t *testing.T) {
	steamResponds(t, `{"response":{}}`)
	if _, err := SteamLevel("76561198056395137"); !errors.Is(err, ErrHidden) {
		t.Errorf("err = %v, want ErrHidden", err)
	}
}
