package main

import (
	"fmt"
	"net/http"
)

type ArmazenamentoJogador interface {
	ObterPontuacaoJogador(nome string) int
}

type ServidorJogador struct {
	armazenamento ArmazenamentoJogador
}

func (s *ServidorJogador) ServerHttp(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogaores/"):]
	fmt.Fprint(w, s.armazenamento.ObterPontuacaoJogador(jogador))
}

func ObterPontuacaoJogador(nome string) string {
	if nome == "Maria" {
		return "20"
	}

	if nome == "Pedro" {
		return "10"
	}

	return ""
}
