package spotifyParser

import (
	"bytes"
	"fmt"
	spotifyAuth "goMusicLinkApi/spotify/auth"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type TrackInfo struct {
	Name   string
	Album  string
	Artist string
}

func ParseSpotifyUrl(url string) {
	fmt.Println("url: ", url)
	trackId := strings.Split(url, "track/")[1]
	getTrackInfo(trackId)

}

// http client to make requests
// TODO: evaluate the feasability to use same client to make all requests to spotify
func client() *http.Client {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	return &client
}

func getTrackInfo(trackId string) {
	url := "https://api.spotify.com/v1/tracks/" + trackId

	token := spotifyAuth.GetToken().AccessToken

	client := *client()

	request, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// The response returns a json containing multiple objects.
	// The objective is to return the most important peaces of information to search on other platforms
	// see details at: https://developer.spotify.com/documentation/web-api/reference/#/operations/get-track
	fmt.Println(string(body))
}

func SearchOnSpotify(info TrackInfo) {
	url := "https://api.spotify.com/v1/search?"

	token := spotifyAuth.GetToken().AccessToken

	client := *client()

	queryPrams := "query=tack%3A" + info.Name + "+artist%3A" + info.Artist + "&type=track&offset=0&limit=20"

	request, err := http.NewRequest("GET", url+queryPrams, bytes.NewBuffer(nil))
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))
}
