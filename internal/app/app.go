package app

import (
	"ctserver/internal/config"
	"ctserver/internal/router"
	"log"
	"net/http"
)

func StartServer() {
	c, err := config.New() // Initialize the config
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	r := router.New(c)
	srv := &http.Server{
		Addr:    ":" + c.Port,
		Handler: r.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
