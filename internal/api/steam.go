package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func ResolveUrl(url string) (string, error) {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	splitUrl := strings.Split(url, "/")
	vanityurl := splitUrl[len(splitUrl)-1]

	if len(splitUrl) < 3 {
		return resolveVanityUrl(vanityurl)
	}

	// TODO: CHECK FOR STEAMID64

	if splitUrl[len(splitUrl)-2] == "id" {
		return resolveVanityUrl(vanityurl)
	}

	if splitUrl[len(splitUrl)-2] == "profiles" {
		return vanityurl, nil
	}

	return "", fmt.Errorf("ResolveUrl: unable to parse url")
}

func resolveVanityUrl(vanityUrl string) (string, error) {
	requestUrl := fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s", os.Getenv("STEAMAPIKEY"), vanityUrl)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return "", fmt.Errorf("ResolveVanityUrl: %w", err)
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var parsed struct {
		Response struct {
			Steamid string `json:"steamid"`
			Success int    `json:"success"`
			Message string `json:"message"`
		} `json:"response"`
	}

	err = json.Unmarshal(bodyBytes, &parsed)
	if err != nil {
		return "", fmt.Errorf("ResolveVanityUrl: %w", err)
	}

	if parsed.Response.Success != 1 {
		return "", fmt.Errorf("ResolveVanityUrl: unable to resolve vanity url")
	}

	return parsed.Response.Steamid, nil
}
