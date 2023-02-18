package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpotifyToken(t *testing.T) {
	sptfToken, err := NewSpotifyToken("tokenkey", "tokenType", 600)
	assert.Nil(t, err)
	assert.NotNil(t, sptfToken)
	assert.NotEmpty(t, sptfToken.AccessToken)
	assert.Equal(t, "tokenkey", sptfToken.AccessToken)
	assert.Equal(t, "tokenType", sptfToken.TokenType)
	assert.Equal(t, uint(600), sptfToken.ExpiresIn)
}

func TestAccessTokenIsRequired(t *testing.T) {
	token, err := NewSpotifyToken("", "type", 600)
	assert.Nil(t, token)
	assert.Equal(t, ErrAccessTokenIsRequired, err)
}

func TestTokenTypeIsRequired(t *testing.T) {
	token, err := NewSpotifyToken("at", "", 600)
	assert.Nil(t, token)
	assert.Equal(t, ErrTokenTypeIsRequired, err)
}

func TestTokenIsExpired(t *testing.T) {
	token, err := NewSpotifyToken("at", "tt", 0)
	assert.Nil(t, token)
	assert.Equal(t, ErrTokenExpired, err)
}
