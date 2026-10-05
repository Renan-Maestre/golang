package main

import (
	"fmt"
	"net/http"
)

func init() {
	http.Handle(
		"/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// time.Sleep(10 * time.Second)
			fmt.Fprintf(w, "Hello, Wold")
		}),
	)

	go func() {
		if err := http.ListenAndServe(":5000", nil); err != nil {
			panic(err)
		}
	}()
}
