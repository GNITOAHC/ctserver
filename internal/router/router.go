package router

import (
	"ctserver/internal/config"
	"ctserver/internal/database"
	"net/http"
)

type Router struct {
	db *database.Database
}

func New(c *config.Config) *Router {
	return &Router{
		db: database.New(c.DatabaseURI),
	}
}

func (r *Router) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive!"))
	})

	return mux
}
