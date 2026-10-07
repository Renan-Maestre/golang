package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"p2/api"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to execute code", "error", err)
		os.Exit(1)
	}

	slog.Info("all systems offline")
}

func run() error {
	apikey := os.Getenv("OMDB_KEY")
	handler := api.Newhandler(apikey)

	s := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
