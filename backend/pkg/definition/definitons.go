package definition

import (
	"fmt"
	"time"
)

var (
	playlistLimit = 50
)

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type Playlist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlaylistResponse struct {
	Items []Playlist `json:"items"`
}

type Song struct {
	Name     string  `json:"name"`
	Artist   string  `json:"artist"`
	Duration float64 `json:"duration"`
}

type UserPlaylist struct {
	Songs []Song
}

type Token struct {
	AccessToken  string    `bson:"access_token"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiryTime   time.Time `bson:"expiry_time,omitempty"`
}

func HandleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func GetPlaylistLimit() int {
	return playlistLimit
}
