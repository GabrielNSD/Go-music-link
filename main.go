package main

import (
	"fmt"

	spotifyParser "goMusicLinkApi/spotify/parser"
)

func main() {
	fmt.Println("Hello music")
	// testInfo := spotifyParser.TrackInfo{
	// 	Name:   "cuff it",
	// 	Album:  "reinascence",
	// 	Artist: "beyonce",
	// }
	spotifyParser.SearchOnSpotify(spotifyParser.ParseSpotifyUrl("http://open.spotify.com/track/6rqhFgbbKwnb9MLmUQDhG6"))
}
