package playing

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
)

func GetCurrentTrack(client *spotify.Client) {
	track, err := client.PlayerCurrentlyPlaying(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if track.Playing {
		fmt.Printf("You are currently playing: %s by %s\n", track.Item.Name, track.Item.Artists[0].Name)
	} else {
		fmt.Println("You are not currently playing anything.")
	}
}
