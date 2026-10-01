package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/horario", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		fmt.Fprintln(w, now)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/user/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "id")
				fmt.Fprintln(w, "rota printa ususrio ", id)
			})
		})
		r.Route("/v2", func(r chi.Router) {
			r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "rota prinat test")
			})
		})

		r.With(middleware.RealIP).
			Get("/realip", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "ip")
			})

		r.Group(func(r chi.Router) {
			r.Use(middleware.BasicAuth("", map[string]string{
				"admin": "admin",
			}))

			r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "pong")
			})
		})
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
