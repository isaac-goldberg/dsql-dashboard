package utils

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/isaac-goldberg/dsql-dashboard/globals"
)

func GetAccessTokens(code string) (string, string, error) {
	var redirectUri string = os.Getenv("HOSTNAME") + "/auth"

	data := url.Values{}
	data.Set("client_id", globals.DISCORD_CLIENT_ID)
	data.Set("client_secret", os.Getenv("DISCORD_CLIENT_SECRET"))
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectUri)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, globals.DiscordTokenUrl, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	json := DecodeJSON(res)

	// if request returned an error
	if json["access_token"] == nil {
		return "", "", errors.New(json["message"].(string))
	}

	return json["access_token"].(string), json["refresh_token"].(string), nil
}

func GetUser(accessToken string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, globals.DiscordUserUrl, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	json := DecodeJSON(res)

	// if request returned an error
	if json["id"] == nil {
		return nil, errors.New(json["message"].(string))
	}

	return json, nil
}

func UseRefreshToken(refreshToken string) (string, string, error) {
	data := url.Values{}
	data.Set("client_id", globals.DISCORD_CLIENT_ID)
	data.Set("client_secret", os.Getenv("DISCORD_CLIENT_SECRET"))
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, globals.DiscordTokenUrl, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)

	if err != nil {
		return "", "", err
	}

	json := DecodeJSON(res)

	if json["error"] != nil {
		return "", "", errors.New(json["error"].(string))
	}

	accessToken := json["access_token"].(string)
	newRefreshToken := json["refresh_token"].(string)

	return accessToken, newRefreshToken, nil
}
