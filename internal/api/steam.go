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
)

type Summary struct {
	Steam64id string
	Avatar    string
	Name      string
	IsPublic  bool
}

func steamGet(method string, params url.Values, out any) error {
	params.Set("key", os.Getenv("STEAMAPIKEY"))

	log.Printf("GET steam %s", method)
	resp, err := steamClient.Get("https://api.steampowered.com/" + method + "/?" + params.Encode())
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

func UserSummary(steam64id string) (Summary, error) {
	var body struct {
		Response struct {
			Players []struct {
				SteamID                  string `json:"steamid"`
				CommunityVisibilityState int    `json:"communityvisibilitystate"`
				PersonaName              string `json:"personaname"`
				AvatarFull               string `json:"avatarfull"`
			} `json:"players"`
		} `json:"response"`
	}
	if err := steamGet("ISteamUser/GetPlayerSummaries/v2", url.Values{"steamids": {steam64id}}, &body); err != nil {
		return Summary{}, fmt.Errorf("UserSummary: %w", err)
	}
	if len(body.Response.Players) == 0 {
		return Summary{}, fmt.Errorf("UserSummary: no player found for %s", steam64id)
	}

	player := body.Response.Players[0]
	return Summary{
		Steam64id: player.SteamID,
		Avatar:    player.AvatarFull,
		Name:      player.PersonaName,
		IsPublic:  player.CommunityVisibilityState == 3,
	}, nil
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
	if err := steamGet("IPlayerService/GetRecentlyPlayedGames/v1", url.Values{"steamid": {steam64id}}, &body); err != nil {
		return 0, fmt.Errorf("UserPlaytime: %w", err)
	}

	for _, game := range body.Response.Games {
		if game.AppID == 730 {
			return int32(game.PlaytimeForever / 60), nil
		}
	}
	return 0, fmt.Errorf("UserPlaytime: no cs playtime found for %s", steam64id)
}
