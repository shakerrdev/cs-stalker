package routes

import (
	"testing"
	"time"
)

func TestThousands(t *testing.T) {
	for n, want := range map[int]string{0: "0", 999: "999", 1000: "1,000", 26437: "26,437", 1234567: "1,234,567", -4500: "-4,500"} {
		if got := thousands(n); got != want {
			t.Errorf("thousands(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestPremier(t *testing.T) {
	tests := []struct {
		rating int
		want   *premierRating
	}{
		{0, nil},
		{850, &premierRating{Low: "850", Tier: 1}},
		{4999, &premierRating{High: "4", Low: ",999", Tier: 1}},
		{5000, &premierRating{High: "5", Low: ",000", Tier: 2}},
		{18234, &premierRating{High: "18", Low: ",234", Tier: 4}},
		{26437, &premierRating{High: "26", Low: ",437", Tier: 6}},
		{30000, &premierRating{High: "30", Low: ",000", Tier: 7}},
		{38012, &premierRating{High: "38", Low: ",012", Tier: 7}},
	}
	for _, tt := range tests {
		got := premier(tt.rating)
		if (got == nil) != (tt.want == nil) || (got != nil && *got != *tt.want) {
			t.Errorf("premier(%d) = %+v, want %+v", tt.rating, got, tt.want)
		}
	}
}

func TestNumberFormats(t *testing.T) {
	tests := []struct{ got, want string }{
		{signed(1.5699999999999998), "+1.57"},
		{signed(-0.4), "-0.40"},
		{rating100(0.0065), "+0.65"},
		{rating100(-0.0318), "-3.18"},
		{points100(0.104), "10.4"},
		{percent(28.3531), "28%"},
		{fraction(0.7037), "70%"},
		{decimal(545.4882, 0), "545"},
		{decimal(2.9374, 1), "2.9"},
		{hours(0), "0 h"},
		{hours(45), "under 1 h"},
		{hours(192645), "3,210 h"},
		{bar(-3), "0.0"},
		{bar(90.1475), "90.1"},
		{bar(130), "100.0"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("got %q, want %q", tt.got, tt.want)
		}
	}
}

func TestAgo(t *testing.T) {
	fixed := time.Date(2026, 10, 12, 20, 0, 0, 0, time.UTC)
	now = func() time.Time { return fixed }
	t.Cleanup(func() { now = time.Now })

	tests := map[time.Duration]string{
		3 * time.Hour:        "today",
		30 * time.Hour:       "yesterday",
		5 * 24 * time.Hour:   "5 days ago",
		20 * 24 * time.Hour:  "2 weeks ago",
		45 * 24 * time.Hour:  "6 weeks ago",
		95 * 24 * time.Hour:  "3 months ago",
		400 * 24 * time.Hour: "13 months ago",
		800 * 24 * time.Hour: "2 years ago",
	}
	for d, want := range tests {
		if got := ago(fixed.Add(-d)); got != want {
			t.Errorf("ago(-%v) = %q, want %q", d, got, want)
		}
	}
}

func TestNames(t *testing.T) {
	eleven, twelve := 11, 12
	tests := []struct{ got, want string }{
		{mapName("de_mirage"), "Mirage"},
		{mapName("de_dust2"), "Dust2"},
		{mapName("cs_italy"), "Italy"},
		{mapName("rush_001"), "Rush 001"},
		{skillGroupShort(1), "S1"},
		{skillGroupShort(15), "LE"},
		{skillGroupShort(18), "GE"},
		{skillGroupShort(0), ""},
		{skillGroupName(18), "The Global Elite"},
		{source("matchmaking", &eleven), "Premier"},
		{source("matchmaking", &twelve), "Matchmaking"},
		{source("matchmaking", nil), "Matchmaking"},
		{source("faceit", nil), "FACEIT"},
		{source("manual_rush", nil), "Manual rush"},
		{outcome("win"), "W"},
		{outcome("loss"), "L"},
		{outcome("tie"), "T"},
		{flag("de"), "🇩🇪"},
		{flag("GB"), "🇬🇧"},
		{flag(""), ""},
		{flag("D1"), ""},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("got %q, want %q", tt.got, tt.want)
		}
	}
}

func TestPresence(t *testing.T) {
	tests := []struct {
		state int
		game  string
		want  presenceState
	}{
		{0, "", presenceState{"Offline", "offline"}},
		{1, "", presenceState{"Online", "online"}},
		{3, "", presenceState{"Away", "online"}},
		{1, "Counter-Strike 2", presenceState{"Playing Counter-Strike 2", "ingame"}},
		{42, "", presenceState{"Offline", "offline"}},
	}
	for _, tt := range tests {
		if got := presence(tt.state, tt.game); got != tt.want {
			t.Errorf("presence(%d, %q) = %+v, want %+v", tt.state, tt.game, got, tt.want)
		}
	}
}

func TestFaceitTier(t *testing.T) {
	for level, want := range map[int]int{1: 1, 2: 2, 3: 2, 4: 3, 7: 3, 8: 4, 9: 4, 10: 5} {
		if got := faceitTier(level); got != want {
			t.Errorf("faceitTier(%d) = %d, want %d", level, got, want)
		}
	}
}

func TestDeref(t *testing.T) {
	five, half := 5, 0.5
	if deref(&five) != 5 || deref(nil) != 0 || derefFloat(&half) != 0.5 || derefFloat(nil) != 0 {
		t.Error("deref")
	}
}
