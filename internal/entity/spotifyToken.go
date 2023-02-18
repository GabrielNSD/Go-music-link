package entity

import (
	"errors"
)

var (
	ErrAccessTokenIsRequired = errors.New("access token is required")
	ErrTokenTypeIsRequired   = errors.New("token type is required")
	ErrTokenExpired          = errors.New("token is expired")
)

type SpotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

func NewSpotifyToken(accessToken, tokenType string, expiresIn uint) (*SpotifyToken, error) {
	token := &SpotifyToken{
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpiresIn:   expiresIn,
	}
	err := token.Validate()
	if err != nil {
		return nil, err
	}

	return token, nil
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
