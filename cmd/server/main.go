package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/GabrielNSD/Go-music-link-api/configs"
	"github.com/GabrielNSD/Go-music-link-api/internal/pkg/spotify"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is the home"))
}

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)

	spotifyToken, err := spotify.NewToken(config.SpotifyClientID, config.SpotifyClientSecret)
	if err != nil {
		log.Println("Spotify token error", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.WithValue("spotifyAccessToken", spotifyToken.AccessToken))
	r.Use(middleware.WithValue("spotifyTokenExpiresIn", spotifyToken.ExpiresIn))
	r.Use(middleware.WithValue("spotifyTokenType", spotifyToken.TokenType))
	r.Use(middleware.WithValue("spotifyTokenExpirationDate", spotifyToken.ExpirationDate))
	r.Use(SpotifyAuthMiddleware)

	r.Use(middleware.Logger)

	r.Get("/", homeHandler)

	http.ListenAndServe(":8001", r)
}

func SpotifyAuthMiddleware(next http.Handler) http.Handler {
	config, _ := configs.LoadConfig(".")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expirationDate := r.Context().Value("spotifyTokenExpirationDate").(time.Time)
		log.Println("exp", expirationDate)
		if expirationDate.After(time.Now()) {
			log.Println("Serving same token")
			next.ServeHTTP(w, r)
		} else {
			log.Println("requesting new token")
			newToken, err := spotify.NewToken(config.SpotifyClientID, config.SpotifyClientSecret)
			if err != nil {
				log.Println("Error generating new Spotify Token")
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, "spotifyAccessToken", newToken.AccessToken)
			ctx = context.WithValue(ctx, "spotifyTokenExpirationDate", newToken.ExpirationDate)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
