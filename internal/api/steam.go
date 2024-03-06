package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func ResolveVanityUrl(url string) (string, error) {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	splitUrl := strings.Split(url, "/")
	vanityurl := splitUrl[len(splitUrl)-1]

	requestUrl := fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s", os.Getenv("STEAMAPIKEY"), vanityurl)
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
		return parsed.Response.Message, fmt.Errorf("ResolveVanityUrl: unable to resolve vanity url")
	}

	return parsed.Response.Steamid, nil
}
