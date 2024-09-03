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

type Url struct {
	Username    string // username
	Url         string // url: Url to shorten
	Path        string // path: Shortened path of the url
	Desc        string // description: Short description of the url
	Ancestor    string // ancestor_id
	ExpireAfter time.Duration
	Default     bool // If true, it's a universal shortened url
}

func (h *Helper) ShortenUrl(u Url) (string, error) {
	if u.Default {
		// Create a default shortened url, require source url, optional path and expireafter
		// Return the shortened url (without base url, just start with _)
		// log.Print(u.Url, u.Path, u.ExpireAfter)
		return h.db.NewUrlDefault(database.Data{
			Content:  u.Url, // url to shorten
			Path:     u.Path,
			Duration: u.ExpireAfter,
		})
	}
	return "", nil
}
