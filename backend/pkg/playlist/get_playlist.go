package playlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func get_playlist(accessToken string) ([]Playlist, error) {
	var user_playlists []Playlist
	for offset := 0; offset < playlistLimit; offset++ {
		cur_playlist, err := get_individual_playlist(accessToken, offset)
		handle_error(err)

		user_playlists = append(user_playlists, cur_playlist...)
		handle_error(err)

		if len(user_playlists) < playlistLimit {
			break
		}

		offset += playlistLimit
	}
	return user_playlists, nil
}

func get_individual_playlist(accessToken string, offset int) ([]Playlist, error) {
	client := &http.Client{}
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?offset=%d&limit=%d", offset, playlistLimit)

	req, err := http.NewRequest("GET", endpoint, nil)
	handle_error(err)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var playlistResponse PlaylistResponse
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		return nil, err
	}

	return playlistResponse.Items, nil
}
