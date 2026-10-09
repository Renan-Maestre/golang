package main

import (
	"log/slog"
	"net/http"
	"time"

	"crudMemoryDB/internal/api"
	"crudMemoryDB/internal/store"
	"crudMemoryDB/internal/user"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		return
	}

	slog.Info("all systems offline")
}

func run() error {
	repository := store.NewUserMemoryRepository()
	service := user.NewService(repository)
	handler := user.NewHandler(service)
	route := api.NewRouter(handler)

	s := http.Server{
		Addr:         ":8080",
		Handler:      route,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
