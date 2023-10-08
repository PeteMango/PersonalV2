package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/PeteMango/website-v2/pkg/auth"
	"github.com/PeteMango/website-v2/pkg/definition"
	"github.com/gin-gonic/gin"
)

var (
	clientID     = "ba45f5ea0ffd4ce59ea1a63dc43baf0d"
	clientSecret = "6c6a3b06bccf449a89fe23e3a5495b80"
	redirectURI  = "http://localhost:8080/callback"
	refreshToken string
)

func main() {
	r := gin.Default()
	var user_playlists []definition.Playlist

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Backend!")
	})

	r.GET("/login", func(c *gin.Context) {
		spotifyAuthURL := "https://accounts.spotify.com/authorize?" + url.Values{
			"response_type": {"code"},
			"client_id":     {clientID},
			"scope":         {"user-read-private user-read-email user-library-read playlist-read-private"},
			"redirect_uri":  {redirectURI},
		}.Encode()
		c.Redirect(http.StatusFound, spotifyAuthURL)
	})

	r.GET("/callback", func(c *gin.Context) {
		code := c.DefaultQuery("code", "")

		accessToken, newRefreshToken, err := auth.ExchangeAccessToken(code)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error exchanging code for access token")
			log.Printf("Error: %v", err)
			return
		}
		refreshToken = newRefreshToken

		playlists, err := get_playlist(accessToken)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching playlists")
			log.Printf("Error: %v", err)
			return
		}
		c.JSON(http.StatusOK, playlists)

		user_playlists = playlists
		for _, playlist := range user_playlists {
			fmt.Printf("Playlist Name: %s\n", playlist.Name)
			fmt.Printf("Playlist ID: %s\n", playlist.ID)
		}

	})

	r.GET("/playlist", func(c *gin.Context) {
		playlistName := c.DefaultQuery("name", "")

		if playlistName == "" {
			c.String(http.StatusBadRequest, "Playlist name is required")
			return
		}

		accessToken, newRefreshToken, err := auth.RefreshAccessToken(refreshToken)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error exchanging code for access token")
			log.Printf("Error: %v", err)
			return
		}
		refreshToken = newRefreshToken

		songs, err := find_playlist(accessToken, user_playlists, playlistName) // Pass playlists here
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching playlist songs")
			fmt.Printf("Error fetching playlist songs: %v\n", err)
			return
		}

		c.JSON(http.StatusOK, songs)

		for _, song := range songs {
			fmt.Printf("%s - %s - %f\n", song.Name, song.Artist, song.Duration)
		}

	})

	r.Run(":8080")
}

func get_playlist(accessToken string) ([]definition.Playlist, error) {
	var user_playlists []definition.Playlist
	for offset := 0; offset < definition.GetPlaylistLimit(); offset++ {
		cur_playlist, err := get_individual_playlist(accessToken, offset)
		definition.HandleError(err)

		user_playlists = append(user_playlists, cur_playlist...)
		definition.HandleError(err)

		if len(user_playlists) < definition.GetPlaylistLimit() {
			break
		}

		offset += definition.GetPlaylistLimit()
	}
	return user_playlists, nil
}

func get_individual_playlist(accessToken string, offset int) ([]definition.Playlist, error) {
	client := &http.Client{}
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?offset=%d&limit=%d", offset, definition.GetPlaylistLimit())

	req, err := http.NewRequest("GET", endpoint, nil)
	definition.HandleError(err)

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

	var playlistResponse definition.PlaylistResponse
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		return nil, err
	}

	return playlistResponse.Items, nil
}

func find_playlist(accessToken string, playlists []definition.Playlist, playlistName string) ([]definition.Song, error) {
	var cur_playlist *definition.Playlist

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

func get_songs(accessToken string, playlistID string) ([]definition.Song, error) {
	var allSongs []definition.Song
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

		var songs []definition.Song
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

			songs = append(songs, definition.Song{Name: name, Artist: artist, Duration: (duration / 1000.0)})
		}

		allSongs = append(allSongs, songs...)

		if next := responseMap["next"]; next == nil || next.(string) == "" {
			break
		}
	}

	return allSongs, nil
}
