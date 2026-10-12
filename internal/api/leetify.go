package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var LeetifyBaseURL = "https://api-public.cs-prod.leetify.com"

type Leetify struct {
	PrivacyMode     string            `json:"privacy_mode"`
	Winrate         float64           `json:"winrate"`
	TotalMatches    int               `json:"total_matches"`
	FirstMatchDate  *time.Time        `json:"first_match_date"`
	Name            string            `json:"name"`
	Bans            []LeetifyBan      `json:"bans"`
	Steam64ID       string            `json:"steam64_id"`
	Ranks           LeetifyRanks      `json:"ranks"`
	Rating          LeetifyRating     `json:"rating"`
	Stats           LeetifyStats      `json:"stats"`
	RecentMatches   []LeetifyMatch    `json:"recent_matches"`
	RecentTeammates []LeetifyTeammate `json:"recent_teammates"`
}

type LeetifyBan struct {
	Platform    string    `json:"platform"`
	Nickname    string    `json:"platform_nickname"`
	BannedSince time.Time `json:"banned_since"`
}

type LeetifyRanks struct {
	Leetify     *float64         `json:"leetify"`
	Premier     *int             `json:"premier"`
	Faceit      *int             `json:"faceit"`
	FaceitElo   *int             `json:"faceit_elo"`
	Wingman     *int             `json:"wingman"`
	Competitive []LeetifyMapRank `json:"competitive"`
}

type LeetifyMapRank struct {
	Map  string `json:"map_name"`
	Rank int    `json:"rank"`
}

type LeetifyRating struct {
	Aim         float64 `json:"aim"`
	Positioning float64 `json:"positioning"`
	Utility     float64 `json:"utility"`
	Clutch      float64 `json:"clutch"`
	Opening     float64 `json:"opening"`
	CTLeetify   float64 `json:"ct_leetify"`
	TLeetify    float64 `json:"t_leetify"`
}

type LeetifyStats struct {
	AccuracyEnemySpotted           float64 `json:"accuracy_enemy_spotted"`
	AccuracyHead                   float64 `json:"accuracy_head"`
	CounterStrafingGoodShots       float64 `json:"counter_strafing_good_shots_ratio"`
	CTOpeningDuelSuccess           float64 `json:"ct_opening_duel_success_percentage"`
	TOpeningDuelSuccess            float64 `json:"t_opening_duel_success_percentage"`
	FlashbangHitFoeAvgDuration     float64 `json:"flashbang_hit_foe_avg_duration"`
	FlashbangHitFoePerFlashbang    float64 `json:"flashbang_hit_foe_per_flashbang"`
	HEFoesDamageAvg                float64 `json:"he_foes_damage_avg"`
	Preaim                         float64 `json:"preaim"`
	ReactionTimeMS                 float64 `json:"reaction_time_ms"`
	SprayAccuracy                  float64 `json:"spray_accuracy"`
	TradedDeathsSuccess            float64 `json:"traded_deaths_success_percentage"`
	TradeKillsSuccess              float64 `json:"trade_kills_success_percentage"`
	UtilityOnDeathAvg              float64 `json:"utility_on_death_avg"`
	FlashbangHitFriendPerFlashbang float64 `json:"flashbang_hit_friend_per_flashbang"`
}

type LeetifyMatch struct {
	ID             string    `json:"id"`
	FinishedAt     time.Time `json:"finished_at"`
	DataSource     string    `json:"data_source"`
	Outcome        string    `json:"outcome"`
	Rank           int       `json:"rank"`
	RankType       *int      `json:"rank_type"`
	Map            string    `json:"map_name"`
	LeetifyRating  float64   `json:"leetify_rating"`
	Score          [2]int    `json:"score"`
	ReactionTimeMS float64   `json:"reaction_time_ms"`
	AccuracyHead   float64   `json:"accuracy_head"`
}

type LeetifyTeammate struct {
	Steam64ID       string `json:"steam64_id"`
	MatchesTogether int    `json:"recent_matches_count"`
}

func LeetifyProfile(steam64id string) (Leetify, error) {
	endpoint := LeetifyBaseURL + "/v3/profile?" + url.Values{"steam64_id": {steam64id}}.Encode()
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return Leetify{}, fmt.Errorf("LeetifyProfile: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if key := os.Getenv("LEETIFY_API_KEY"); key != "" {
		req.Header.Set("_leetify_key", key)
	}

	log.Printf("GET leetify profile")
	res, err := leetifyClient.Do(req)
	if err != nil {
		return Leetify{}, fmt.Errorf("LeetifyProfile: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return Leetify{}, fmt.Errorf("LeetifyProfile: %s: %w", steam64id, ErrNotFound)
	default:
		return Leetify{}, fmt.Errorf("LeetifyProfile: %s", res.Status)
	}

	var profile Leetify
	if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
		return Leetify{}, fmt.Errorf("LeetifyProfile: decode: %w", err)
	}
	return profile, nil
}
