package spotifyAuth

import (
	"context"
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

	pgx "github.com/jackc/pgx/v5"

	// temporary use, maybe it is better to set env variables in a docker container
	// this package reads the dotenv and sets its variables to environment
	"github.com/joho/godotenv"
)

type SpotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

type DatabaseToken struct {
	AccessToken  string
	TokenType    string
	Scope        string
	Expiration   time.Time
	RefreshToken string
}

func getTokenFromDB() *DatabaseToken {
	pgURL := "postgres://useruser:password@localhost:5431/music-link"
	db, err := pgx.Connect(context.Background(), pgURL)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer db.Close(context.Background())

	serviceName := "spotify"
	var token DatabaseToken
	err = db.QueryRow(context.Background(), "select access_token, token_type, scope, expiration, refresh_token from tokens where service_name=$1", serviceName).Scan(&token.AccessToken, &token.TokenType, &token.Scope, &token.Expiration, &token.RefreshToken)
	if err != nil {
		return nil
	}

	now := time.Now()
	fmt.Println(now)
	fmt.Println(token.Expiration)

	if now.After(token.Expiration) {
		return nil
	}

	fmt.Println(token)
	return &token
}

func writeTokenToDB(token *SpotifyToken) {
	pgURL := "postgres://useruser:password@localhost:5431/music-link"
	db, err := pgx.Connect(context.Background(), pgURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	expirationTime := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))

	commandTag, err := db.Exec(context.Background(), "insert into tokens(service_name, access_token, token_type, scope, expiration, refresh_token) values($1, $2, $3, $4, $5, $6)", "spotify", token.AccessToken, token.TokenType, "", expirationTime, "")
	if err != nil {
		fmt.Println(err)
	}

	// if the insertion was not possible, update existing token data
	if commandTag.RowsAffected() != 1 {
		_, err := db.Exec(context.Background(),
			"update tokens set access_token=$1, token_type=$2, scope=$3, expiration=$4, refresh_token=$5 where service_name=$6", token.AccessToken, token.TokenType, "", expirationTime, "", "spotify")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func requestNewToken() (*SpotifyToken, error) {
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

	// TODO: check if the returned json does contain only an error key
	var token SpotifyToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &token, nil
}

func GetToken() (*SpotifyToken, error) {
	dbToken := getTokenFromDB()
	if dbToken != nil {
		var token SpotifyToken
		token.AccessToken = dbToken.AccessToken
		token.TokenType = dbToken.TokenType
		currentTime := time.Now()
		token.ExpiresIn = uint(currentTime.Sub(dbToken.Expiration).Seconds())
		fmt.Println("returning from DB")
		return &token, nil
	}

	token, err := requestNewToken()
	if err != nil {
		log.Fatal(err)
	}

	writeTokenToDB(token)

	return token, nil
}
