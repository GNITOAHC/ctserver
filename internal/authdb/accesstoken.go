package authdb

import (
	"ctserver/jwt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type AccessToken struct {
	Mail     string
	Username string
	gojwt.RegisteredClaims
}

func (a *AccessToken) ToDomain(m map[string]interface{}) {
	a.Mail = m["Mail"].(string)
	a.Username = m["Username"].(string)
	return
}

func (a *AccessToken) SetClaims(claims gojwt.RegisteredClaims) {
	a.RegisteredClaims = claims
}

// NewAccessToken creates a new access token
func NewAccessToken(mail, username, secret, prefix string, duration time.Duration) (string, error) {
	signed, err := jwt.Sign(&AccessToken{Mail: mail, Username: username}, secret, prefix, duration)
	if err != nil {
		return "", err
	}
	return signed, nil
}
