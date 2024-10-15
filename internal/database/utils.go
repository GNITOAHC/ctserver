package database

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"math/big"
	"strconv"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// generateShortened generates a hashed value for the given url
func generateShortened(url string) string {
	// Generate MD5
	hash := md5.Sum([]byte(url)) // Generate MD5
	hashPrefix := hash[:6]

	num := new(big.Int).SetBytes(hashPrefix) // Convert hash bytes to big int

	// Encode the big int to base62
	shortened := func(num *big.Int) string {
		result := ""
		base := big.NewInt(62)
		zero := big.NewInt(0)

		for num.Cmp(zero) > 0 {
			mod := new(big.Int)
			num.DivMod(num, base, mod)
			result = string(base62Chars[mod.Int64()]) + result
		}

		return result
	}(num)

	return shortened
}

func rehash(url string, existing map[string]bool) (string, error) {
	shortened := generateShortened(url)
	counter := 0
	for {
		counter++
		newUrl := generateShortened(url + strconv.Itoa(counter))
		if _, exists := existing[newUrl]; !exists {
			return shortened, nil
		}
		if counter > 100 {
			break
		}
	}
	return "", errors.New("hash collision")
}

// NullString returns a sql.NullString from the given string
func NullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// DefaultExpiredAt returns the default expired time for the given duration
// If duration is 0, default to 180 days; If duration is -1, no expiration
func ExpiredAt(t time.Duration) sql.NullTime {
	if t == -1 {
		return sql.NullTime{Time: time.Time{}, Valid: false}
	}
	if t == 0 {
		t = time.Hour * 24 * 180
	}
	return sql.NullTime{Time: time.Now().Add(t), Valid: true}
}
