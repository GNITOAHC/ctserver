package database

import (
	"context"
	"errors"
	"time"
)

// NewUrlDefault creates a new shortened URL for the given original url
// custom is the custom shortened path (optional), if set, should start with _
func (db *Database) NewUrlDefault(original, custom string, duration time.Duration) (string, error) {
	// Check required fields
	if original == "" {
		return "", errors.New("content is required")
	}
	if custom == "" {
		custom = "_" + generateShortened(original)
	}
	if custom[0] != '_' {
		custom = "_" + custom
	}
	query := `
		with new_data as (insert into data_t (username, path, type, description, ancestor_id, content, expired_at) values ($1, $2, 'url', $3, $4, $5, $6) returning id),
			 insert_shortened as (insert into shortened_t (data_id, shortened) values ((select id from new_data), $2) returning *)
		select shortened from insert_shortened;
	`
	var resultShortened string
	var d DataTable
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

// NewData creates a column for the given data, require d.Path, d.Content and d.Username
// d.Path should start with /, so is the returned string
func (db *Database) NewData(d DataTable) (string, error) {
	// Check required fields
	if d.Path == "" || d.Content == "" || d.Username == "" {
		return "", errors.New("path, content and username are required")
	}
	query := `
		insert into data_t (username, path, type, description, ancestor_id, content, expired_at) values ($1, $2, $3, $4, $5, $6, $7) returning path;
	`
	var resultPath string
	err := db.pool.QueryRow(context.Background(), query,
		d.Username, d.Path, d.Type, d.Description, NullString(d.AncestorId), d.Content, ExpiredAt(d.Duration),
	).Scan(&resultPath)
	if err != nil {
		return "", err
	}
	return resultPath, nil
}

// Return the original url for the given username and path
func (db *Database) GetData(username, path string) (string, error) {
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

func (db *Database) GetChildren(username, path string) ([]DataTable, error) {
	query := `
        select (id, type, content, description) from data_t where ancestor_id = (select id from data_t where username = $1 and path = $2 and type = 'dir');
    `
	var result []DataTable
	rows, err := db.pool.Query(context.Background(), query, username, path)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var d DataTable
		err := rows.Scan(&d.Id, &d.Type, &d.Content, &d.Description)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}
