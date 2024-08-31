package helper

import "ctserver/internal/database"

type Helper struct {
	db *database.Database
}

func New(dburi string) *Helper {
	return &Helper{db: database.New(dburi)}
}

func (h *Helper) CheckUserExist(mail string) (bool, error) {
	return h.db.CheckUserExist(mail)
}

func (h *Helper) RegisterUser(mail, phone string) error {
	return h.db.InsertUser(mail, phone)
}
