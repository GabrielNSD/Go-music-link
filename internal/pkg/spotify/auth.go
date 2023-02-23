package spotify

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	// "os"
	"errors"
	"strings"
	"time"
	// temporary use, maybe it is better to set env variables in a docker container
	// this package reads the dotenv and sets its variables to environment
	// "github.com/joho/godotenv"
)

type SpotifyToken struct {
	AccessToken    string `json:"access_token"`
	TokenType      string `json:"token_type"`
	ExpiresIn      uint   `json:"expires_in"`
	ExpirationDate time.Time
}

type SpotifyTokenData struct {
	AccessToken    string
	ExpirationDate time.Time
}

var (
	ErrAccessTokenIsRequired = errors.New("access token is required")
	ErrTokenTypeIsRequired   = errors.New("token type is required")
	ErrTokenExpired          = errors.New("token is expired")
)

func NewToken(clientId, clientSecret string) (*SpotifyToken, error) {
	//func GetToken() (*SpotifyToken, error) {
	// err := godotenv.Load("local.env")
	// if err != nil {
	// 	log.Fatalf("Error loading env file: %s", err)
	// 	return nil, err
	// }
	// clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	// clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

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

	t := time.Now()
	token.ExpirationDate = t.Add(time.Second * time.Duration(token.ExpiresIn))

	log.Println("New token", token)

	return &token, nil
}

func (t *SpotifyToken) GetToken() string {
	if t.ExpirationDate.After(time.Now()) {
		return t.AccessToken
	}
	return "a"
	// t = NewToken()
	// return t.AccessToken
}

func (t *SpotifyToken) Validate() error {
	if t.AccessToken == "" {
		return ErrAccessTokenIsRequired
	}
	if t.TokenType == "" {
		return ErrTokenTypeIsRequired
	}
	if t.ExpiresIn == 0 {
		return ErrTokenExpired
	}

	return nil
}
