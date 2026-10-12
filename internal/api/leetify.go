package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Profile struct {
	RecentGameRatings struct {
		Aim                 float64 `json:"aim"`
		LeetifyRatingRounds int     `json:"leetifyRatingRounds"`
		Positioning         float64 `json:"positioning"`
		Utility             float64 `json:"utility"`
		GamesPlayed         int     `json:"gamesPlayed"`
		Clutch              float64 `json:"clutch"`
		CtLeetify           float64 `json:"ctLeetify"`
		Leetify             float64 `json:"leetify"`
		Opening             float64 `json:"opening"`
		TLeetify            float64 `json:"tLeetify"`
	} `json:"recentGameRatings"`
	Teammates []struct {
		IsBanned                 bool    `json:"isBanned"`
		MatchesPlayedTogether    int     `json:"matchesPlayedTogether"`
		ProfileUserLeetifyRating float64 `json:"profileUserLeetifyRating"`
		Rank                     struct {
			Type       string `json:"type"`
			DataSource string `json:"dataSource"`
			SkillLevel int    `json:"skillLevel"`
		} `json:"rank"`
		Steam64ID             string  `json:"steam64Id"`
		SteamAvatarURL        string  `json:"steamAvatarUrl"`
		SteamNickname         string  `json:"steamNickname"`
		TeammateLeetifyRating float64 `json:"teammateLeetifyRating"`
		WinRateTogether       float64 `json:"winRateTogether"`
	} `json:"teammates"`
	Meta struct {
		Name           string   `json:"name"`
		Steam64ID      string   `json:"steam64Id"`
		SteamAvatarURL string   `json:"steamAvatarUrl"`
		FaceitNickname string   `json:"faceitNickname"`
		PlatformBans   []string `json:"platformBans"`
	} `json:"meta"`
}

func LeetifyProfile(steam64id string) (Profile, error) {
	endpoints := []string{
		"https://api.leetify.com/api/profile/id/" + url.PathEscape(steam64id),
		"https://api.leetify.com/api/profile/" + url.PathEscape(steam64id),
	}

	for _, endpoint := range endpoints {
		profile, found, err := fetchLeetifyProfile(endpoint)
		if err != nil {
			return Profile{}, fmt.Errorf("LeetifyProfile: %w", err)
		}
		if found {
			return profile, nil
		}
	}

	return Profile{}, fmt.Errorf("LeetifyProfile: no profile found for %s", steam64id)
}

func fetchLeetifyProfile(endpoint string) (Profile, bool, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return Profile{}, false, err
	}
	req.Header.Set("Accept", "application/json")

	log.Printf("GET %s", endpoint)
	res, err := leetifyClient.Do(req)
	if err != nil {
		return Profile{}, false, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return Profile{}, false, nil
	default:
		return Profile{}, false, fmt.Errorf("GET %s: %s", endpoint, res.Status)
	}

	var profile Profile
	if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
		return Profile{}, false, fmt.Errorf("decode %s: %w", endpoint, err)
	}
	return profile, true, nil
}
