package main

import (
	"log/slog"
	"net/http"
	"time"

	"projeto1/api"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to exuce code", "error", err)
		return
		// os.Exit(1)
	}
	slog.Info("all systems offline")
}

func run() error {
	db := make(map[string]string)
	handler := api.NewHandler(db)

	s := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
