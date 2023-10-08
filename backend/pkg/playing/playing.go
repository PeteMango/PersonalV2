package playing

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func GetCurrentlyPlayingSong(accessToken string) (string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", "", errors.New(string(bodyBytes))
	}

	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return "", "", err
	}

	itemMap, ok := responseMap["item"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("Unable to decode currently playing song data")
	}

	name, ok := itemMap["name"].(string)
	if !ok {
		return "", "", errors.New("Unable to extract song name from the data")
	}

	artists, ok := itemMap["artists"].([]interface{})
	if !ok || len(artists) == 0 {
		return "", "", errors.New("Unable to extract artist data from the response")
	}

	artistMap, ok := artists[0].(map[string]interface{})
	if !ok {
		return "", "", errors.New("Unable to decode artist data")
	}

	artistName, ok := artistMap["name"].(string)
	if !ok {
		return "", "", errors.New("Unable to extract artist name from the data")
	}

	return name, artistName, nil
}
