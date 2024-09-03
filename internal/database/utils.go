package database

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"time"
)

// generateShortened generates a hashed value for the given url
func generateShortened(url string) string {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8] // Return the first 8 characters of the hashed value
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
