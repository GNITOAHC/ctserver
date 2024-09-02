package authdb

import (
	"ctserver/jwt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	Mail         string // User's mail
	SessionToken string // Is the session id
	gojwt.RegisteredClaims
}

func (r *RefreshToken) ToDomain(m map[string]interface{}) {
	r.SessionToken = m["SessionToken"].(string)
	r.Mail = m["Mail"].(string)
	return
}

func (r *RefreshToken) SetClaims(claims gojwt.RegisteredClaims) {
	r.RegisteredClaims = claims
}

// Sign signs the token and retuns the signed token with given prefix
func Sign(token *RefreshToken, secret, prefix string, duration time.Duration) (string, error) {
	signed, err := jwt.Sign(token, secret, prefix, duration)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func Decode(token string, secret, prefix string) (*RefreshToken, error) {
	decoded, err := jwt.Parse(token, secret, prefix)
	if err != nil {
		return nil, err
	}
	if isExpired, err := jwt.IsExpired(decoded); isExpired || err != nil {
		return nil, err
	}
	var refreshToken RefreshToken
	refreshToken.ToDomain(decoded)
	return &refreshToken, nil
}
