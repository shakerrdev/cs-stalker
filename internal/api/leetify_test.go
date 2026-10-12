package api

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"
)

func leetifyResponds(t *testing.T, status int, body []byte) *http.Request {
	t.Helper()
	var got http.Request
	useLeetifyServer(t, func(w http.ResponseWriter, r *http.Request) {
		got = *r.Clone(r.Context())
		w.WriteHeader(status)
		w.Write(body)
	})
	return &got
}

func TestLeetifyProfile(t *testing.T) {
	fixture, err := os.ReadFile("testdata/leetify_profile.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("LEETIFY_API_KEY", "")
	req := leetifyResponds(t, http.StatusOK, fixture)

	profile, err := LeetifyProfile("76561198056395137")
	if err != nil {
		t.Fatal(err)
	}

	if req.URL.Path != "/v3/profile" || req.URL.Query().Get("steam64_id") != "76561198056395137" {
		t.Errorf("request = %s", req.URL)
	}
	if req.Header.Get("_leetify_key") != "" {
		t.Error("sent an api key header without LEETIFY_API_KEY")
	}
	if profile.Winrate != 0.5512 || profile.TotalMatches != 812 || profile.PrivacyMode != "public" {
		t.Errorf("profile = %+v", profile)
	}
	if profile.FirstMatchDate == nil || !profile.FirstMatchDate.Equal(time.Date(2021, 3, 14, 18, 2, 11, 0, time.UTC)) {
		t.Errorf("first match date = %v", profile.FirstMatchDate)
	}
	ranks := profile.Ranks
	if *ranks.Premier != 18234 || *ranks.Faceit != 7 || *ranks.FaceitElo != 1650 || *ranks.Wingman != 12 || *ranks.Leetify != 1.5699999999999998 {
		t.Errorf("ranks = %+v", ranks)
	}
	if len(ranks.Competitive) != 2 || ranks.Competitive[0] != (LeetifyMapRank{Map: "de_mirage", Rank: 15}) {
		t.Errorf("map ranks = %+v", ranks.Competitive)
	}
	if profile.Rating.Aim != 74.1475 || profile.Rating.Opening != -0.018 || profile.Rating.TLeetify != 0.0276 {
		t.Errorf("rating = %+v", profile.Rating)
	}
	if profile.Stats.ReactionTimeMS != 545.4882 || profile.Stats.AccuracyHead != 28.3531 {
		t.Errorf("stats = %+v", profile.Stats)
	}
	if len(profile.Bans) != 1 || profile.Bans[0].Platform != "faceit" || profile.Bans[0].Nickname != "tester_fc" {
		t.Errorf("bans = %+v", profile.Bans)
	}

	if len(profile.RecentMatches) != 2 {
		t.Fatalf("recent matches = %d", len(profile.RecentMatches))
	}
	first, second := profile.RecentMatches[0], profile.RecentMatches[1]
	if first.Map != "de_anubis" || first.Outcome != "win" || first.Score != [2]int{13, 8} || first.RankType == nil || *first.RankType != 11 || first.LeetifyRating != 0.048 {
		t.Errorf("first match = %+v", first)
	}
	if second.RankType != nil || second.DataSource != "faceit" {
		t.Errorf("second match = %+v", second)
	}
	if len(profile.RecentTeammates) != 2 || profile.RecentTeammates[0] != (LeetifyTeammate{Steam64ID: "76561198000000001", MatchesTogether: 33}) {
		t.Errorf("teammates = %+v", profile.RecentTeammates)
	}
}

func TestLeetifyProfile_SendsAPIKey(t *testing.T) {
	t.Setenv("LEETIFY_API_KEY", "secret")
	req := leetifyResponds(t, http.StatusOK, []byte(`{}`))

	if _, err := LeetifyProfile("76561198056395137"); err != nil {
		t.Fatal(err)
	}
	if got := req.Header.Values("_leetify_key"); len(got) != 1 || got[0] != "secret" {
		t.Errorf("_leetify_key = %q", got)
	}
}

func TestLeetifyProfile_NotFound(t *testing.T) {
	leetifyResponds(t, http.StatusNotFound, []byte(`{"error":"not found"}`))
	if _, err := LeetifyProfile("76561198056395137"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestLeetifyProfile_Errors(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{"server error", http.StatusInternalServerError, ""},
		{"rate limited", http.StatusTooManyRequests, ""},
		{"malformed json", http.StatusOK, "<html>maintenance</html>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leetifyResponds(t, tt.status, []byte(tt.body))
			_, err := LeetifyProfile("76561198056395137")
			if err == nil || errors.Is(err, ErrNotFound) {
				t.Fatalf("err = %v, want a non-not-found error", err)
			}
		})
	}
}
