package router

import (
	"ctserver/cache"
	"ctserver/internal/authdb"
	"ctserver/internal/config"
	"ctserver/internal/helper"
	"ctserver/mailer"
	"encoding/json"
	"net/http"
	"time"
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
	mux.HandleFunc("POST /shorten-url", r.UniversalShortenUrl)

	// For testing
	mux.HandleFunc("POST /test-delete", r.TestDelete)

	wrapped := use(mux, middleware)

	return wrapped
}

func (rr *Router) UniversalShortenUrl(w http.ResponseWriter, r *http.Request) {
	type ShortenedReq struct {
		Source      string        `json:"source"`
		CustomPath  string        `json:"customPath,omitempty"`
		ExpireAfter time.Duration `json:"expireAfter,omitempty"`
	}
	var req ShortenedReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Source == "" {
		http.Error(w, "Source is required", http.StatusBadRequest)
		return
	}

	shortened, err := rr.helper.ShortenUrl(req.Source, req.CustomPath, req.ExpireAfter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rr.config.BaseURL + "/" + shortened)
}
