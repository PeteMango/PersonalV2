package definition

import "fmt"

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

func HandleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func GetPlaylistLimit() int {
	return playlistLimit
}
