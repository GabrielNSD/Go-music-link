package main

import (
	"fmt"

	spotifyAuth "goMusicLinkApi/spotify"
)

func main() {
	fmt.Println("Hello music")
	spotifyAuth.GetToken() // testing package
}
