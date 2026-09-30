package main

import (
	"log"
	"net/http"
)

func main() {
	servidor := &ServidorJogador{}
	if err := http.ListenAndServe(":5000", servidor); err != nil {
		log.Fatalf("não foi possível escutar na porta 5000 %v", err)
	}
}
