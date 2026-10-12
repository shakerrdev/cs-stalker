package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shakerrdev/cs-stalker/internal/api"
)

const (
	playerID = "76561198056395137"
	mateID   = "76561198000000001"
)

type upstream struct {
	players       map[string]map[string]any
	leetify       []byte
	leetifyStatus int
	ownedGames    []map[string]any
	bans          map[string]any
}

func newUpstream(t *testing.T) *upstream {
	t.Helper()
	leetify, err := os.ReadFile("testdata/leetify_profile.json")
	if err != nil {
		t.Fatal(err)
	}
	return &upstream{
		players: map[string]map[string]any{
			playerID: {"steamid": playerID, "personaname": "tester", "communityvisibilitystate": 3, "personastate": 1, "avatarfull": "https://avatars.example/full.jpg", "profileurl": "https://steamcommunity.com/id/tester/", "timecreated": 1300000000, "loccountrycode": "DE"},
			mateID:   {"steamid": mateID, "personaname": "duo partner", "communityvisibilitystate": 3, "avatarmedium": "https://avatars.example/mate.jpg"},
		},
		leetify:       leetify,
		leetifyStatus: http.StatusOK,
		ownedGames:    []map[string]any{{"appid": 730, "playtime_forever": 192600, "playtime_2weeks": 840, "rtime_last_played": 1760190000}},
		bans:          map[string]any{"SteamId": playerID, "NumberOfVACBans": 0, "NumberOfGameBans": 0, "EconomyBan": "none"},
	}
}

func (u *upstream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	switch strings.TrimSuffix(r.URL.Path, "/") {
	case "/ISteamUser/ResolveVanityURL/v1":
		if query.Get("vanityurl") == "tester" {
			replyJSON(w, map[string]any{"response": map[string]any{"success": 1, "steamid": playerID}})
			return
		}
		replyJSON(w, map[string]any{"response": map[string]any{"success": 42, "message": "No match"}})
	case "/ISteamUser/GetPlayerSummaries/v2":
		players := []map[string]any{}
		for _, id := range strings.Split(query.Get("steamids"), ",") {
			if player, ok := u.players[id]; ok {
				players = append(players, player)
			}
		}
		replyJSON(w, map[string]any{"response": map[string]any{"players": players}})
	case "/ISteamUser/GetPlayerBans/v1":
		replyJSON(w, map[string]any{"players": []any{u.bans}})
	case "/IPlayerService/GetSteamLevel/v1":
		replyJSON(w, map[string]any{"response": map[string]any{"player_level": 42}})
	case "/IPlayerService/GetOwnedGames/v1":
		replyJSON(w, map[string]any{"response": map[string]any{"games": u.ownedGames}})
	case "/v3/profile":
		w.WriteHeader(u.leetifyStatus)
		w.Write(u.leetify)
	default:
		http.NotFound(w, r)
	}
}

func replyJSON(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func serveUpstream(t *testing.T, u *upstream) {
	t.Helper()
	server := httptest.NewServer(u)
	t.Cleanup(server.Close)

	steam, leetify := api.SteamBaseURL, api.LeetifyBaseURL
	api.SteamBaseURL, api.LeetifyBaseURL = server.URL, server.URL
	t.Cleanup(func() { api.SteamBaseURL, api.LeetifyBaseURL = steam, leetify })

	previous := now
	now = func() time.Time { return time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { now = previous })
}

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

func assertContains(t *testing.T, rec *httptest.ResponseRecorder, wants ...string) {
	t.Helper()
	body := rec.Body.String()
	for _, want := range wants {
		if !strings.Contains(body, want) {
			t.Errorf("body does not contain %q", want)
		}
	}
}

func TestPages(t *testing.T) {
	router := newTestRouter(t)
	tests := []struct {
		target   string
		contains string
	}{
		{"/", `action="/lookup"`},
		{"/banana", `id="minion"`},
		{"/static/index.css", "box-sizing"},
		{"/static/fonts/saira.woff2", "wOF2"},
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

func TestHealthz(t *testing.T) {
	rec := get(newTestRouter(t), "/healthz")
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Errorf("GET /healthz = %d %q", rec.Code, rec.Body.String())
	}
}

func TestLookupRedirectsToPlayer(t *testing.T) {
	serveUpstream(t, newUpstream(t))
	router := newTestRouter(t)
	for _, query := range []string{playerID, "tester", "https://steamcommunity.com/id/tester/", "https://steamcommunity.com/profiles/" + playerID} {
		rec := get(router, "/lookup?q="+url.QueryEscape(query))
		if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/player/"+playerID {
			t.Errorf("lookup %q = %d %q", query, rec.Code, rec.Header().Get("Location"))
		}
	}
}

func TestLookupWithoutQueryGoesHome(t *testing.T) {
	rec := get(newTestRouter(t), "/lookup?q=+")
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/" {
		t.Errorf("lookup = %d %q", rec.Code, rec.Header().Get("Location"))
	}
}

func TestLookupShowsErrors(t *testing.T) {
	serveUpstream(t, newUpstream(t))
	router := newTestRouter(t)
	for _, query := range []string{"https://example.com/not-a-profile", "nobody-here"} {
		rec := get(router, "/lookup?q="+url.QueryEscape(query))
		if rec.Code != http.StatusNotFound {
			t.Errorf("lookup %q: status = %d", query, rec.Code)
		}
		assertContains(t, rec, "No Steam profile matches that", `value="`+query+`"`)
	}
}

func TestBananaLookup(t *testing.T) {
	rec := get(newTestRouter(t), "/lookup?q=Banana")
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/banana" {
		t.Errorf("lookup = %d %q", rec.Code, rec.Header().Get("Location"))
	}
}

func TestPlayerPage(t *testing.T) {
	u := newUpstream(t)
	u.bans["NumberOfVACBans"] = 1
	u.bans["DaysSinceLastBan"] = 340
	serveUpstream(t, u)

	rec := get(newTestRouter(t), "/player/"+playerID)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body.String())
	}
	assertContains(t, rec,
		"<title>tester – cs-stalker</title>",
		"Steam level 42",
		"On Steam since Mar 2011",
		`<span class="hi">18</span><span class="lo">,234</span>`,
		"1 VAC ban",
		"Last ban 340 days ago",
		"FACEIT ban since May 2024 as tester_fc",
		"3,210 h",
		"Level 7",
		"1,650 Elo",
		"55%",
		"Anubis",
		"11 Oct",
		"&#43;4.80",
		`href="/player/`+mateID+`"`,
		"duo partner",
		"33 matches",
		"76561198000000002",
		`title="Legendary Eagle"`,
	)
	if strings.Contains(rec.Body.String(), "Overpass") {
		t.Error("unranked map is listed")
	}
}

func TestPlayerPageWithoutLeetify(t *testing.T) {
	u := newUpstream(t)
	u.leetifyStatus = http.StatusNotFound
	u.leetify = []byte(`{"error":"not found"}`)
	u.ownedGames = nil
	serveUpstream(t, u)

	rec := get(newTestRouter(t), "/player/"+playerID)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	assertContains(t, rec, "Leetify has no data on this player", "Their Steam game details are private", "No VAC, game or trade bans on record")
}

func TestPlayerPageErrors(t *testing.T) {
	serveUpstream(t, newUpstream(t))
	router := newTestRouter(t)
	tests := []struct {
		target string
		want   string
	}{
		{"/player/not-an-id", "Paste the profile link instead"},
		{"/player/76561198999999999", "No Steam account has that ID"},
	}
	for _, tt := range tests {
		rec := get(router, tt.target)
		if rec.Code != http.StatusNotFound {
			t.Errorf("GET %s: status = %d", tt.target, rec.Code)
		}
		assertContains(t, rec, tt.want)
	}
}
