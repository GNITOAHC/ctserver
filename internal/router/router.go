package router

import (
	"ctserver/internal/config"
	"net/http"
)

type Router struct {
}

func New(c *config.Config) *Router {
	return &Router{}
}

func (r *Router) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive!"))
	})

	return mux
}
