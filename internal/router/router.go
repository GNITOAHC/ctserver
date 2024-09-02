package router

import (
	"ctserver/cache"
	"ctserver/internal/authdb"
	"ctserver/internal/config"
	"ctserver/internal/helper"
	"ctserver/mailer"
	"net/http"
)

type Router struct {
	helper *helper.Helper
	mailer mailer.Mailer
	config *config.Config
	cache  *cache.Cache[string, string]
	authdb *authdb.RefreshDB
}

func New(c *config.Config) *Router {
	return &Router{
		helper: helper.New(c.DatabaseURI),
		mailer: mailer.New(c.SMTPFrom, c.SMTPPass, c.SMTPHost, c.SMTPPort),
		config: c,
		cache:  cache.New[string, string](),
		authdb: authdb.New(c.AuthDBURI, c.AuthDBName, c.AuthDBCollection, c.JWTSecret),
	}
}

func (r *Router) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive!"))
	})

	mux.HandleFunc("POST /register", r.Register)
	mux.HandleFunc("POST /register/verify", r.RegVerify)
	mux.HandleFunc("POST /login", r.Login)
	mux.HandleFunc("POST /login/verify", r.LoginVerify)

	return mux
}
