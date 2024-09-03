package database

import (
	"context"
	"errors"
	"time"
)

type Data struct {
	Username     string
	Id           string
	Path         string // Path of the data, should start with /
	Type         string // "file", "url", "dir" or "text"
	Content      string // Original URL for type "url", source url for type "file", or content of dir or text
	Description  string // Short description of the data
	AncestorId   string
	DescendantId []string
	CreatedAt    time.Time
	ExpiredAt    time.Time

	Duration time.Duration
}

// NewUrlDefault creates a new shortened URL for the given data, require d.Content as original url
// d.Path is the shortened path (optional), if set, should start with _
func (db *Database) NewUrlDefault(d Data) (string, error) {
	// Check required fields
	if d.Content == "" {
		return "", errors.New("content is required")
	}
	if d.Path == "" {
		d.Path = "_" + generateShortened(d.Content)
	}
	query := `
		with new_data as (insert into data_t (username, path, type, description, ancestor_id, content, expired_at) values ($1, $2, 'url', $3, $4, $5, $6) returning id),
			 insert_shortened as (insert into shortened_t (data_id, shortened) values ((select id from new_data), $2) returning *)
		select shortened from insert_shortened;
	`
	var resultShortened string
	err := db.pool.QueryRow(context.Background(), query, d.Username, d.Path, d.Description, NullString(d.AncestorId), d.Content, ExpiredAt(d.Duration)).Scan(&resultShortened)
	if err != nil {
		return "", err
	}
	return resultShortened, nil
}

// GetUrlDefault get the original URL for the given shortened URL
func (db *Database) GetUrlDefault(shortened string) (string, error) {
	query := `
		select d.content 
		from data_t d
		join shortened_t s on d.id = s.data_id
		where s.shortened = $1;
	`
	var originalUrl string
	err := db.pool.QueryRow(context.Background(), query, shortened).Scan(&originalUrl)
	if err != nil {
		return "", err
	}
	return originalUrl, nil
}

// NewUrlCustom creates a new shortened URL for the given data, require d.Path, d.Content and d.Username
// d.Path should start with /, so is the returned string
func (db *Database) NewUrlCustom(d Data) (string, error) {
	// Check required fields
	if d.Path == "" || d.Content == "" || d.Username == "" {
		return "", errors.New("path, content and username are required")
	}
	query := `
		insert into data_t (username, path, type, description, ancestor_id, content, expired_at) values ($1, $2, 'url', $3, $4, $5, $6) returning path;
	`
	var resultPath string
	err := db.pool.QueryRow(context.Background(), query, d.Username, d.Path, d.Description, NullString(d.AncestorId), d.Content, ExpiredAt(d.Duration)).Scan(&resultPath)
	if err != nil {
		return "", err
	}
	return resultPath, nil
}

// Return the original url for the given username and path
func (db *Database) GetUrlCustom(username, path string) (string, error) {
	query := `
		select d.content
		from data_t d
		where d.username = $1 and d.path = $2;
	`
	var content string
	err := db.pool.QueryRow(context.Background(), query, username, path).Scan(&content)
	if err != nil {
		return "", err
	}
	return content, nil
}
