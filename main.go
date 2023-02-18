package main

import (
	"fmt"

	"github.com/GabrielNSD/Go-music-link-api/spotify/parser"
)

func main() {
	fmt.Println("Hello music")
	testInfo := spotifyParser.TrackInfo{
		Name:   "cuff it",
		Album:  "renaissance",
		Artist: "beyonce",
	}
	// testPink := spotifyParser.TrackInfo{
	// 	Name:   "speak to me - 2011 Remastered Version",
	// 	Album:  "the dark side of the moon",
	// 	Artist: "pink floyd",
	// }
	// spotifyParser.SearchOnSpotify(testPink)
	spotifyParser.SearchOnSpotify(testInfo)
	// The response for this track is not being returned
	// spotifyParser.SearchOnSpotify(spotifyParser.ParseSpotifyUrl("http://open.spotify.com/track/6rqhFgbbKwnb9MLmUQDhG6"))
}
