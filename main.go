package main

import (
	"fmt"

	spotifyParser "goMusicLinkApi/spotify/parser"
)

func main() {
	fmt.Println("Hello music")
	spotifyParser.ParseSpotifyUrl("http://open.spotify.com/track/6rqhFgbbKwnb9MLmUQDhG6")
	testInfo := spotifyParser.TrackInfo{
		Name:   "cozy",
		Album:  "reinascence",
		Artist: "beyonce",
	}
	spotifyParser.SearchOnSpotify(testInfo)
}
