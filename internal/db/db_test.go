package db

import (
	"os"
	"slices"
	"testing"

	"github.com/shakerrdev/cs-stalker/internal/db/models"
	"github.com/shakerrdev/cs-stalker/internal/db/models/faceit"
	"github.com/shakerrdev/cs-stalker/internal/db/models/leetify"
	"github.com/shakerrdev/cs-stalker/internal/db/models/steam"
	"gorm.io/gorm"
)

func testTx(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	if err := Connect(dsn); err != nil {
		t.Fatal(err)
	}
	tx := DB.Begin()
	t.Cleanup(func() { tx.Rollback() })
	return tx
}

func newPlayer(steam64 string, club leetify.LeetifyClub) models.Player {
	return models.Player{
		SteamID64:     steam64,
		Name:          "player " + steam64,
		Bans:          []string{"vac"},
		SteamProfile:  steam.SteamProfile{Steam64ID: steam64, SteamName: "steam " + steam64},
		FaceitProfile: faceit.FaceitProfile{FaceitID: "faceit-" + steam64, Nickname: "faceit " + steam64, ELO: 2000},
		LeetifyProfile: leetify.LeetifyProfile{
			LeetifyUserID: "leetify-" + steam64,
			Steam64ID:     steam64,
			Aim:           71.5,
			Club:          club,
			Teammates: []leetify.LeetifyTeammate{{
				Steam64Id:             "76561198000000001",
				MatchesPlayedTogether: 12,
				Rank:                  leetify.LeetifyRank{Type: "premier", SkillLevel: 15000},
			}},
			Highlights:    []leetify.LeetifyHighlight{{Description: "ace", RoundNumber: 7}},
			PersonalBests: []leetify.LeetifyPersonalBest{{SkillID: "aim", Value: "99"}},
		},
	}
}

func TestPlayerRoundTrip(t *testing.T) {
	tx := testTx(t)

	player := newPlayer("76561198056395137", leetify.LeetifyClub{ClubID: "club-1", Name: "Club", Tag: "CLB"})
	if err := tx.Create(&player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}

	var got models.Player
	err := tx.
		Preload("SteamProfile").
		Preload("FaceitProfile").
		Preload("LeetifyProfile.Club").
		Preload("LeetifyProfile.Teammates.Rank").
		Preload("LeetifyProfile.Highlights").
		Preload("LeetifyProfile.PersonalBests").
		First(&got, player.ID).Error
	if err != nil {
		t.Fatalf("load player: %v", err)
	}

	if !slices.Equal(got.Bans, []string{"vac"}) {
		t.Errorf("bans = %v, want [vac]", got.Bans)
	}
	if got.SteamProfile.SteamName != player.SteamProfile.SteamName {
		t.Errorf("steam name = %q", got.SteamProfile.SteamName)
	}
	if got.LeetifyProfile.Club.Name != "Club" {
		t.Errorf("club = %+v", got.LeetifyProfile.Club)
	}
	if len(got.LeetifyProfile.Teammates) != 1 || got.LeetifyProfile.Teammates[0].Rank.SkillLevel != 15000 {
		t.Errorf("teammates = %+v", got.LeetifyProfile.Teammates)
	}
	if len(got.LeetifyProfile.Highlights) != 1 || len(got.LeetifyProfile.PersonalBests) != 1 {
		t.Fatalf("highlights = %d, personal bests = %d", len(got.LeetifyProfile.Highlights), len(got.LeetifyProfile.PersonalBests))
	}
	if got.LeetifyProfile.Highlights[0].CreatedAt.IsZero() {
		t.Error("highlight created_at was not set")
	}
}

func TestPlayersCanShareAClub(t *testing.T) {
	tx := testTx(t)

	var ids []uint
	for _, steam64 := range []string{"76561198000000010", "76561198000000011"} {
		player := newPlayer(steam64, leetify.LeetifyClub{ClubID: "shared", Name: "Shared"})
		player.FaceitProfile = faceit.FaceitProfile{}
		if err := tx.Create(&player).Error; err != nil {
			t.Fatalf("create player %s: %v", steam64, err)
		}
		ids = append(ids, player.ID)
	}

	for _, id := range ids {
		var got models.Player
		if err := tx.Preload("LeetifyProfile.Club").First(&got, id).Error; err != nil {
			t.Fatalf("load player %d: %v", id, err)
		}
		if got.LeetifyProfile.Club.Name != "Shared" {
			t.Errorf("player %d lost its club: %+v", id, got.LeetifyProfile.Club)
		}
	}
}

func TestMatchRoundTrip(t *testing.T) {
	tx := testTx(t)

	player := newPlayer("76561198056395137", leetify.LeetifyClub{})
	if err := tx.Create(&player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}

	match := models.Match{
		Map:        "de_mirage",
		ScoreTeam1: 13,
		ScoreTeam2: 9,
		LeetifyMatch: leetify.LeetifyMatch{
			LeetifyID: "match-1",
			MapName:   "de_mirage",
			Details:   leetify.LeetifyMatchDetails{Tickrate: 64, Ticks: 180000},
			PlayerStats: []leetify.LeetifyPlayerStats{
				{Steam64Id: player.SteamID64, Kills: 25, Deaths: 14, Adr: 92.4},
			},
			Agents:   []leetify.LeetifyAgent{{Steam64Id: player.SteamID64, TeamNumber: 2}},
			Parties:  []leetify.LeetifyParty{{Steam64Id: player.SteamID64, Party: 1}},
			Profiles: []leetify.LeetifyProfile{player.LeetifyProfile},
		},
	}
	if err := tx.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}

	var got models.Match
	err := tx.
		Preload("LeetifyMatch.Details").
		Preload("LeetifyMatch.PlayerStats").
		Preload("LeetifyMatch.Agents").
		Preload("LeetifyMatch.Parties").
		Preload("LeetifyMatch.Profiles").
		First(&got, match.ID).Error
	if err != nil {
		t.Fatalf("load match: %v", err)
	}

	if got.UUID == "" {
		t.Error("uuid was not generated")
	}
	lm := got.LeetifyMatch
	if lm.Details.Ticks != 180000 || len(lm.PlayerStats) != 1 || len(lm.Agents) != 1 || len(lm.Parties) != 1 {
		t.Errorf("leetify match = %+v", lm)
	}
	if len(lm.Profiles) != 1 || lm.Profiles[0].Steam64ID != player.SteamID64 {
		t.Errorf("profiles = %+v", lm.Profiles)
	}

	var profile leetify.LeetifyProfile
	if err := tx.Preload("Matches").First(&profile, player.LeetifyProfile.ID).Error; err != nil {
		t.Fatalf("load profile: %v", err)
	}
	if len(profile.Matches) != 1 || profile.Matches[0].LeetifyID != "match-1" {
		t.Errorf("profile matches = %+v", profile.Matches)
	}
}
