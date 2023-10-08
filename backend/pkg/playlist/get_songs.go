package playlist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func find_playlist(accessToken string, playlists []Playlist, playlistName string) ([]Song, error) {
	var cur_playlist *Playlist

	for _, playlist := range playlists {
		if playlist.Name == playlistName {
			cur_playlist = &playlist
			break
		}
	}

	if cur_playlist == nil {
		return nil, fmt.Errorf("Playlist not found")
	}

	return get_songs(accessToken, cur_playlist.ID)
}

func get_songs(accessToken string, playlistID string) ([]Song, error) {
	var allSongs []Song
	client := &http.Client{}

	limit := 100 // Spotify's maximum limit
	for offset := 0; ; offset += limit {
		apiURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?limit=%d&offset=%d", playlistID, limit, offset)

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var responseMap map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
			return nil, err
		}

		items, ok := responseMap["items"]
		if !ok || items == nil {
			return nil, fmt.Errorf("No items found in the playlist")
		}

		tracks, ok := items.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Couldn't type assert tracks to []interface{}")
		}

		var songs []Song
		for _, item := range tracks {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			trackMap, ok := itemMap["track"].(map[string]interface{})
			if !ok {
				continue
			}

			name, ok := trackMap["name"].(string)
			if !ok {
				continue
			}

			artistsSlice, ok := trackMap["artists"].([]interface{})
			if !ok || len(artistsSlice) == 0 {
				continue
			}

			artistMap, ok := artistsSlice[0].(map[string]interface{}) // First artist
			if !ok {
				continue
			}

			artist, ok := artistMap["name"].(string)
			if !ok {
				continue
			}

			duration, ok := trackMap["duration_ms"].(float64)
			if !ok {
				continue
			}

			songs = append(songs, Song{Name: name, Artist: artist, Duration: (duration / 1000.0)})
		}

		allSongs = append(allSongs, songs...)

		if next := responseMap["next"]; next == nil || next.(string) == "" {
			break
		}
	}

	return allSongs, nil
}

// func get_songs(accessToken string, playlistID string) ([]Song, error) {
// 	client := &http.Client{}
// 	apiURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

// 	req, err := http.NewRequest("GET", apiURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Authorization", "Bearer "+accessToken)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// Decode the response and print it for debugging
// 	var responseMap map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
// 		return nil, err
// 	}
// 	fmt.Printf("Response: %+v\n", responseMap)

// 	items, ok := responseMap["items"]
// 	if !ok || items == nil {
// 		fmt.Println("No items in response or items is nil.")
// 		return nil, fmt.Errorf("No items found in the playlist")
// 	}

// 	tracks, ok := items.([]interface{})
// 	if !ok {
// 		fmt.Println("Items could not be type asserted to []interface{}.")
// 		return nil, fmt.Errorf("Couldn't type assert tracks to []interface{}")
// 	}

// 	var songs []Song
// 	for _, item := range tracks {
// 		itemMap, ok := item.(map[string]interface{})
// 		if !ok {
// 			fmt.Println("Item could not be type asserted to map[string]interface{}.")
// 			continue
// 		}
// 		trackMap, ok := itemMap["track"].(map[string]interface{})
// 		if !ok {
// 			fmt.Println("Track could not be type asserted to map[string]interface{}.")
// 			continue
// 		}
// 		name, ok := trackMap["name"].(string)
// 		if !ok {
// 			fmt.Println("Name could not be type asserted to string.")
// 			continue
// 		}

// 		artistsSlice, ok := trackMap["artists"].([]interface{})
// 		if !ok || len(artistsSlice) == 0 {
// 			fmt.Println("Artists could not be type asserted to []interface{} or no artists found.")
// 			continue
// 		}
// 		artistMap, ok := artistsSlice[0].(map[string]interface{}) // Assuming you want the first artist
// 		if !ok {
// 			fmt.Println("First artist could not be type asserted to map[string]interface{}.")
// 			continue
// 		}
// 		artist, ok := artistMap["name"].(string)
// 		if !ok {
// 			fmt.Println("Artist name could not be type asserted to string.")
// 			continue
// 		}

// 		duration, ok := trackMap["duration_ms"].(float64)
// 		if !ok {
// 			fmt.Println("Song duration could not be type asserted to integer.")
// 			continue
// 		}

// 		songs = append(songs, Song{Name: name, Artist: artist, Duration: (duration / 1000.0)})
// 	}

// 	return songs, nil
// }
