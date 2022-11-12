package spotifyParser

import (
	"bytes"
	"encoding/json"
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

func ParseSpotifyUrl(url string) TrackInfo {
	fmt.Println("url: ", url)
	trackId := strings.Split(url, "track/")[1]
	return getTrackInfo(trackId)

}

type ParsedAlbum struct {
	Name string `json:"name"`
}

type ParsedArtist struct {
	Href string `json:"href"`
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Uri  string `json:"uri"`
}

type ParsedTrack struct {
	Album   ParsedAlbum    `json:"album"`
	Artists []ParsedArtist `json:"artists"`
	Name    string         `json:"name"`
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

func getTrackInfo(trackId string) TrackInfo {
	url := "https://api.spotify.com/v1/tracks/" + trackId

	token, err := spotifyAuth.GetToken()

	if err != nil {
		log.Fatalln(err)
	}

	client := *client()

	request, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))

	request.Header.Add("Authorization", "Bearer "+token.AccessToken)
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

	var parsedTrack ParsedTrack
	err = json.Unmarshal(body, &parsedTrack)
	if err != nil {
		log.Fatalln(err)
	}

	var result = TrackInfo{
		Name:   parsedTrack.Name,
		Album:  parsedTrack.Album.Name,
		Artist: parsedTrack.Artists[0].Name,
	}

	return result
}

func SearchOnSpotify(info TrackInfo) {
	url := "https://api.spotify.com/v1/search?"

	token, err := spotifyAuth.GetToken()

	if err != nil {
		log.Fatalln(err)
	}

	client := *client()

	q := strings.ReplaceAll("q="+"track:"+info.Name+"%20artist:"+info.Artist, " ", "%20")
	typeParam := "&type=track"

	queryParams := q + typeParam

	fmt.Println("Param", queryParams)

	request, err := http.NewRequest("GET", url+queryParams, bytes.NewBuffer(nil))
	request.Header.Add("Authorization", "Bearer "+token.AccessToken)
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

	fmt.Println("search result")
	fmt.Println(string(body))
}
