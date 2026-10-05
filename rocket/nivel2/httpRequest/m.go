package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	time.Sleep(1 * time.Second)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:5000", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("ACCEPT", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
