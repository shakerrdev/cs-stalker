package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
		fmt.Sprintf("https://api.leetify.com/api/profile/id/%s", steam64id),
		fmt.Sprintf("https://api.leetify.com/api/profile/%s", steam64id),
	}

	for _, url := range endpoints {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return Profile{}, fmt.Errorf("LeetifyProfile: %w", err)
		}

		req.Header.Add("Accept", "application/json, text/plain, */*")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("LEETIFY_AUTH_TOKEN")))

		log.Printf("GET %s", url)
		res, err := leetifyClient.Do(req)
		if err != nil {
			return Profile{}, fmt.Errorf("LeetifyProfile: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode == 200 {
			body, _ := io.ReadAll(res.Body)

			var profile Profile
			json.Unmarshal(body, &profile)

			return profile, nil
		}
	}

	return Profile{}, fmt.Errorf("LeetifyProfile: no valid profile found (tried both /id and plain endpoints)")
}
