package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

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
	var user_playlists []Playlist

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

		accessToken, newRefreshToken, err := exchange_accessToken(code)
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

		accessToken, newRefreshToken, err := refresh_accessToken(refreshToken)
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
