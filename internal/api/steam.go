package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	ErrHidden   = errors.New("hidden by privacy settings")
)

type Summary struct {
	Steam64id    string
	Name         string
	Avatar       string
	AvatarMedium string
	ProfileURL   string
	IsPublic     bool
	PersonaState int
	Game         string
	Created      time.Time
	CountryCode  string
}

var SteamBaseURL = "https://api.steampowered.com"

func steamGet(method string, params url.Values, out any) error {
	params.Set("key", os.Getenv("STEAMAPIKEY"))

	log.Printf("GET steam %s", method)
	resp, err := steamClient.Get(SteamBaseURL + "/" + method + "/?" + params.Encode())
	if err != nil {
		// url.Error repeats the request URL, which contains the API key.
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return fmt.Errorf("GET %s: %w", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", method, resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s: %w", method, err)
	}
	return nil
}

func resolveVanityUrl(vanityUrl string) (string, error) {
	var body struct {
		Response struct {
			SteamID string `json:"steamid"`
			Success int    `json:"success"`
		} `json:"response"`
	}
	if err := steamGet("ISteamUser/ResolveVanityURL/v1", url.Values{"vanityurl": {vanityUrl}}, &body); err != nil {
		return "", fmt.Errorf("ResolveVanityUrl: %w", err)
	}
	if body.Response.Success != 1 {
		return "", fmt.Errorf("ResolveVanityUrl: no steam id found for %q", vanityUrl)
	}
	return body.Response.SteamID, nil
}

var steam64Pattern = regexp.MustCompile(`^7656\d{13}$`)

func IsSteam64(id string) bool {
	return steam64Pattern.MatchString(id)
}

func ResolveUrl(input string) (string, error) {
	input = strings.TrimSpace(input)
	switch {
	case input == "":
		return "", fmt.Errorf("ResolveUrl: empty input")
	case steam64Pattern.MatchString(input):
		return input, nil
	case !strings.Contains(input, "/"):
		return resolveVanityUrl(input)
	}

	if i := strings.IndexAny(input, "?#"); i >= 0 {
		input = input[:i]
	}
	segments := strings.Split(strings.TrimSuffix(input, "/"), "/")
	if len(segments) >= 2 {
		kind, value := segments[len(segments)-2], segments[len(segments)-1]
		switch {
		case kind == "id" && value != "":
			return resolveVanityUrl(value)
		case kind == "profiles" && steam64Pattern.MatchString(value):
			return value, nil
		}
	}

	return "", fmt.Errorf("ResolveUrl: not a steam profile url: %q", input)
}

func PlayerSummaries(steam64ids []string) (map[string]Summary, error) {
	summaries := make(map[string]Summary, len(steam64ids))
	for start := 0; start < len(steam64ids); start += 100 {
		batch := steam64ids[start:min(start+100, len(steam64ids))]
		var body struct {
			Response struct {
				Players []struct {
					SteamID                  string `json:"steamid"`
					CommunityVisibilityState int    `json:"communityvisibilitystate"`
					PersonaName              string `json:"personaname"`
					ProfileURL               string `json:"profileurl"`
					AvatarMedium             string `json:"avatarmedium"`
					AvatarFull               string `json:"avatarfull"`
					PersonaState             int    `json:"personastate"`
					GameExtraInfo            string `json:"gameextrainfo"`
					TimeCreated              int64  `json:"timecreated"`
					LocCountryCode           string `json:"loccountrycode"`
				} `json:"players"`
			} `json:"response"`
		}
		params := url.Values{"steamids": {strings.Join(batch, ",")}}
		if err := steamGet("ISteamUser/GetPlayerSummaries/v2", params, &body); err != nil {
			return nil, fmt.Errorf("PlayerSummaries: %w", err)
		}

		for _, player := range body.Response.Players {
			summary := Summary{
				Steam64id:    player.SteamID,
				Name:         player.PersonaName,
				Avatar:       player.AvatarFull,
				AvatarMedium: player.AvatarMedium,
				ProfileURL:   player.ProfileURL,
				IsPublic:     player.CommunityVisibilityState == 3,
				PersonaState: player.PersonaState,
				Game:         player.GameExtraInfo,
				CountryCode:  player.LocCountryCode,
			}
			if player.TimeCreated > 0 {
				summary.Created = time.Unix(player.TimeCreated, 0).UTC()
			}
			summaries[player.SteamID] = summary
		}
	}
	return summaries, nil
}

func UserSummary(steam64id string) (Summary, error) {
	summaries, err := PlayerSummaries([]string{steam64id})
	if err != nil {
		return Summary{}, fmt.Errorf("UserSummary: %w", err)
	}
	summary, ok := summaries[steam64id]
	if !ok {
		return Summary{}, fmt.Errorf("UserSummary: %s: %w", steam64id, ErrNotFound)
	}
	return summary, nil
}

func UserPlaytime(steam64id string) (int32, error) {
	var body struct {
		Response struct {
			Games []struct {
				AppID           int `json:"appid"`
				PlaytimeForever int `json:"playtime_forever"`
			} `json:"games"`
		} `json:"response"`
	}
	params := url.Values{"steamid": {steam64id}, "include_played_free_games": {"1"}}
	if err := steamGet("IPlayerService/GetOwnedGames/v1", params, &body); err != nil {
		return 0, fmt.Errorf("UserPlaytime: %w", err)
	}

	for _, game := range body.Response.Games {
		if game.AppID == 730 {
			return int32(game.PlaytimeForever / 60), nil
		}
	}
	return 0, fmt.Errorf("UserPlaytime: no cs playtime found for %s", steam64id)
}

type Bans struct {
	VACBans          int
	GameBans         int
	DaysSinceLastBan int
	CommunityBanned  bool
	TradeBan         string
}

func (b Bans) Any() bool {
	return b.VACBans > 0 || b.GameBans > 0 || b.CommunityBanned || b.TradeBan != ""
}

func PlayerBans(steam64id string) (Bans, error) {
	var body struct {
		Players []struct {
			CommunityBanned  bool
			NumberOfVACBans  int
			NumberOfGameBans int
			DaysSinceLastBan int
			EconomyBan       string
		} `json:"players"`
	}
	if err := steamGet("ISteamUser/GetPlayerBans/v1", url.Values{"steamids": {steam64id}}, &body); err != nil {
		return Bans{}, fmt.Errorf("PlayerBans: %w", err)
	}
	if len(body.Players) == 0 {
		return Bans{}, fmt.Errorf("PlayerBans: %s: %w", steam64id, ErrNotFound)
	}

	player := body.Players[0]
	bans := Bans{
		VACBans:          player.NumberOfVACBans,
		GameBans:         player.NumberOfGameBans,
		DaysSinceLastBan: player.DaysSinceLastBan,
		CommunityBanned:  player.CommunityBanned,
	}
	if player.EconomyBan != "" && player.EconomyBan != "none" {
		bans.TradeBan = player.EconomyBan
	}
	return bans, nil
}

func SteamLevel(steam64id string) (int, error) {
	var body struct {
		Response struct {
			PlayerLevel *int `json:"player_level"`
		} `json:"response"`
	}
	if err := steamGet("IPlayerService/GetSteamLevel/v1", url.Values{"steamid": {steam64id}}, &body); err != nil {
		return 0, fmt.Errorf("SteamLevel: %w", err)
	}
	if body.Response.PlayerLevel == nil {
		return 0, fmt.Errorf("SteamLevel: %s: %w", steam64id, ErrHidden)
	}
	return *body.Response.PlayerLevel, nil
}
