package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"projetoComDb/internal/api"
	"projetoComDb/internal/store"

	"github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to exuce code", "error", err)
		os.Exit(1)
	}
	slog.Info("all systems offline")
}

func run() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "",
		DB:       0,
	})

	store := store.NewStore(rdb)
	handler := api.NewHandler(store)

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
