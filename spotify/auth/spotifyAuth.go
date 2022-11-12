package spotifyAuth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	// temporary use, maybe it is better to set env variables in a docker container
	// this package reads the dotenv and sets its variables to environment
	"github.com/joho/godotenv"
)

type SpotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

func GetToken() (*SpotifyToken, error) {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error loading env file: %s", err)
		return nil, err
	}
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	apiUrl := "https://accounts.spotify.com/api/token"

	encodedSecrets := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))

	formData := url.Values{
		"grant_type": {"client_credentials"},
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", apiUrl, strings.NewReader((formData.Encode())))

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Basic "+encodedSecrets)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var token SpotifyToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	fmt.Printf("%+v\n", token)

	return &token, nil
}
