package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PeteMango/website-v2/pkg/auth"
	"github.com/PeteMango/website-v2/pkg/db"
	"github.com/PeteMango/website-v2/pkg/playing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://localhost:8080/callback"
	refreshToken string
)

func main() {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Add the origin of your React app
	r := gin.Default()
	r.Use(cors.New(config))
	mongoClient := db.InitializeMongo()
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	coll := mongoClient.Database("songs").Collection("previous_plays")

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Backend!")
	})

	r.GET("/login", func(c *gin.Context) {
		spotifyAuthURL := "https://accounts.spotify.com/authorize?" + url.Values{
			"response_type": {"code"},
			"client_id":     {clientID},
			"scope":         {"user-read-playback-state user-read-private user-read-email user-library-read playlist-read-private"},
			"redirect_uri":  {redirectURI},
		}.Encode()
		c.Redirect(http.StatusFound, spotifyAuthURL)
	})

	r.GET("/callback", func(c *gin.Context) {
		tokenColl := mongoClient.Database("songs").Collection("tokens")

		existingToken, expiry, err := db.FetchTokens(tokenColl)
		if err != nil {
			log.Printf("Failed to get tokens: %v", err)
		}

		var accessToken string
		if existingToken != nil && expiry.After(time.Now()) {
			accessToken = existingToken.AccessToken
			refreshToken = existingToken.RefreshToken
		} else {
			code := c.DefaultQuery("code", "")
			accessToken, newRefreshToken, err := auth.ExchangeAccessToken(code)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error exchanging code for access token")
				log.Printf("Error: %v", err)
				return
			}
			refreshToken = newRefreshToken
			fmt.Printf("THE REFRESH TOKEN IS %s\n", refreshToken)

			err = db.UpsertTokens(tokenColl, accessToken, refreshToken, time.Now().Add(1*time.Hour))
			if err != nil {
				log.Printf("Failed to upsert tokens: %v", err)
			}
		}

		songName, songArtist, songLink, err := playing.GetCurrentlyPlayingSong(accessToken)
		if err != nil {
			fmt.Printf("THE ERROR IS: %s", err)
			c.String(http.StatusInternalServerError, "Error fetching currently playing song")
			log.Printf("Error: %v", err)
			return
		}

		result, err := db.InsertSongs(coll, songName, songArtist, songLink)
		if err != nil {
			log.Fatalf("Failed to insert document: %v", err)
		}

		fmt.Println(result)

		c.JSON(http.StatusOK, gin.H{
			"songName":   songName,
			"songArtist": songArtist,
			"songLink":   songLink,
		})
	})

	r.GET("/recentSong", func(c *gin.Context) {
		coll := mongoClient.Database("songs").Collection("previous_plays")
		result, err := db.FetchMostRecentSong(coll)

		if err != nil {
			fmt.Printf("failed to access db\n")
		}

		c.JSON(http.StatusOK, gin.H{
			"songName":   result["song"],
			"songArtist": result["artist"],
			"songLink":   result["link"],
		})
	})
	r.Run(":8080")
}
