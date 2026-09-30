package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestObterJogador(t *testing.T) {
	servidor := &ServidorJogador{}
	t.Run("retorna resultado de Maria", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontucao("Maria")
		respotas := httptest.NewRecorder()

		servidor.ServerHttp(respotas, requisicao)

		recebido := respotas.Body.String()
		esperado := "20"

		verificaCorpoRequisicao(t, recebido, esperado)
	})

	t.Run("retorna resultado de Pedro", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontucao("Pedro")
		resposta := httptest.NewRecorder()

		servidor.ServerHttp(resposta, requisicao)

		recebido := resposta.Body.String()
		esperado := "10"

		verificaCorpoRequisicao(t, recebido, esperado)
	})
}

func novaRequisicaoObterPontucao(nome string) *http.Request {
	requisicao, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return requisicao
}

func verificaCorpoRequisicao(t *testing.T, recebido, esperado string) {
	t.Helper()

	if recebido != esperado {
		t.Errorf("recebido '%s', esperado '%s'", recebido, esperado)
	}
}
