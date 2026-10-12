package routes

import (
	"reflect"
	"testing"
	"time"

	"github.com/shakerrdev/cs-stalker/internal/api"
)

func TestDescribeBans(t *testing.T) {
	faceit := api.LeetifyBan{Platform: "faceit", Nickname: "fc_name", BannedSince: time.Date(2024, 5, 2, 0, 0, 0, 0, time.UTC)}
	tests := []struct {
		name      string
		bans      *api.Bans
		leetify   []api.LeetifyBan
		wantState string
		wantLines []string
	}{
		{"steam unavailable", nil, nil, "unknown", nil},
		{"clean", &api.Bans{}, nil, "clean", nil},
		{"vac and game bans", &api.Bans{VACBans: 2, GameBans: 1, DaysSinceLastBan: 1}, nil, "banned", []string{"2 VAC bans", "1 game ban", "Last ban yesterday"}},
		{"community and trade", &api.Bans{CommunityBanned: true, TradeBan: "probation"}, nil, "banned", []string{"Banned from the Steam Community", "On trade probation"}},
		{"trade ban", &api.Bans{TradeBan: "banned"}, nil, "banned", []string{"Trade banned"}},
		{"faceit only", &api.Bans{}, []api.LeetifyBan{faceit}, "banned", []string{"FACEIT ban since May 2024 as fc_name"}},
		{"faceit while steam unavailable", nil, []api.LeetifyBan{{Platform: "esea", BannedSince: faceit.BannedSince}}, "banned", []string{"ESEA ban since May 2024"}},
	}
	for _, tt := range tests {
		state, lines := describeBans(tt.bans, tt.leetify)
		if state != tt.wantState || !reflect.DeepEqual(lines, tt.wantLines) {
			t.Errorf("%s: got %q %q, want %q %q", tt.name, state, lines, tt.wantState, tt.wantLines)
		}
	}
}

func TestRankedMaps(t *testing.T) {
	got := rankedMaps([]api.LeetifyMapRank{{Map: "de_overpass", Rank: 0}, {Map: "de_nuke", Rank: 9}, {Map: "de_mirage", Rank: 15}, {Map: "de_inferno", Rank: 15}})
	want := []api.LeetifyMapRank{{Map: "de_inferno", Rank: 15}, {Map: "de_mirage", Rank: 15}, {Map: "de_nuke", Rank: 9}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("rankedMaps = %v, want %v", got, want)
	}
}
