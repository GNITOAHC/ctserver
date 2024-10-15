package helper

import (
	"ctserver/internal/database"
	"time"
)

type Helper struct {
	db *database.Database
}

func New(dburi string) *Helper {
	return &Helper{db: database.New(dburi)}
}

func (h *Helper) CheckUserExist(mail string) (bool, error) {
	return h.db.CheckUserExist(mail)
}

func (h *Helper) CheckUsernameExist(username string) (bool, error) {
	return h.db.CheckUsernameExist(username)
}

func (h *Helper) RegisterUser(mail, phone, username string) error {
	return h.db.InsertUser(mail, phone, username)
}

func (h *Helper) GetUsername(mail string) (string, error) {
	return h.db.GetUsername(mail)
}

func (h *Helper) RemoveUser(mail string) error {
	return h.db.RemoveUser(mail)
}

// ShortenUrl creates a new shortened URL for the given original url (with optional custom path and expiration duration)
func (h *Helper) ShortenUrl(url, custom string, duration time.Duration) (string, error) {
	return h.db.NewUrlDefault(url, custom, duration)
}
