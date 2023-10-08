package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://localhost:8080/callback"
	refreshToken string
)

func ExchangeAccessToken(code string) (string, string, error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	response, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", err
	}

	response.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(response)
	if resp.StatusCode != http.StatusOK {
		return "", "", err
	}

	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Received Access Token: ", tokenResponse.AccessToken)
	fmt.Println("Received Refresh Token: ", tokenResponse.RefreshToken)
	return tokenResponse.AccessToken, tokenResponse.RefreshToken, nil
}
func RefreshAccessToken(refreshToken string) (string, string, error) {
	tokenURL := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Updated Access Token: ", tokenResponse.AccessToken)
	fmt.Println("Updated Refresh Token: ", tokenResponse.RefreshToken)
	return tokenResponse.AccessToken, tokenResponse.RefreshToken, nil
}
