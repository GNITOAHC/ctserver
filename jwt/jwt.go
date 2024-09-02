package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims interface {
	ToDomain(m map[string]interface{})
	SetClaims(claims jwt.RegisteredClaims)
	jwt.Claims
}

// Sign sign a token with given content and key. Returns the signed token.
// `content`: should be a struct that implements TokenClaims interface.
func Sign(content TokenClaims, key, prefix string, t time.Duration) (string, error) {
	// Set up claims
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(t)), // 2 day expiration
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	c := content
	c.SetClaims(claims)

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	// Sign the token with the key
	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}
	return prefix + tokenString, nil
}

// Parse function parse the token and return the claims as a map[string]interface{}.
func Parse(tokenString, key, prefix string) (map[string]interface{}, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}

	token, err := jwt.Parse(tokenString[len(prefix):], keyFunc, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}

	// Validate the token and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Convert MapClaims to a map[string]interface{}
		result := make(map[string]interface{})
		for key, value := range claims {
			result[key] = value
		}
		return result, nil
	}

	return nil, errors.New("Invalid token")
}

func IsExpired(claims map[string]interface{}) (bool, error) {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, errors.New("Invalid expiration time")
	}
	return time.Now().Unix() > int64(exp), nil
}
