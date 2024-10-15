package database

import (
	"bytes"
	"time"
)

type RequiredField struct {
	Username   string
	Path       string
	AncestorId string // Empty if root
}

// Create operation's request the data_t
type CreateRequest struct {
	// Required either one of these fields
	File    *bytes.Buffer
	Url     string
	DirName string
	Text    string

	Description string        // Optional
	Duration    time.Duration // Optional
}

// CreateFile creates a new file for the user
func (db *Database) CreateFile(rf RequiredField, cr CreateRequest) (string, error) {
	// Upload file to storage and get the URL
	mockurl := "https://storage.com/file"

	// Insert the file to the database
	d := DataTable{
		Username:    rf.Username,
		Path:        rf.Path,
		Type:        "file",
		Content:     mockurl,
		Description: cr.Description,
		AncestorId:  rf.AncestorId,
		Duration:    cr.Duration,
	}
	path, err := db.NewData(d)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (db *Database) CreateUrl(rf RequiredField, cr CreateRequest) (string, error) {
	d := DataTable{
		Username:    rf.Username,
		Path:        rf.Path,
		Type:        "url",
		Content:     cr.Url,
		Description: cr.Description,
		AncestorId:  rf.AncestorId,
		Duration:    cr.Duration,
	}
	path, err := db.NewData(d)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (db *Database) CreateDir(rf RequiredField, cr CreateRequest) (string, error) {
	d := DataTable{
		Username:    rf.Username,
		Path:        rf.Path,
		Type:        "dir",
		Content:     cr.DirName,
		Description: cr.Description,
		AncestorId:  rf.AncestorId,
		Duration:    cr.Duration,
	}
	path, err := db.NewData(d)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (db *Database) CreateText(cr CreateRequest) error {
	return nil
}
