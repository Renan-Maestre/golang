package main

import "testing"

func TestOla(t *testing.T) {
    verificaMensagemCorreta := func(t testing.TB, resultado, esperado string) {
        t.Helper()
        if resultado != esperado {
            t.Errorf("resultado %q, esperado %q", resultado, esperado)
        }
    }

    t.Run("diz olá para as pessoas", func(t *testing.T) {
        resultado := Ola("Renan", "")
        esperado := "Olá, Renan"
        verificaMensagemCorreta(t, resultado, esperado)
    })

    t.Run("'Mundo' como padrão para string vazia", func(t *testing.T) {
        resultado := Ola("", "")
        esperado := "Olá, Mundo"
        verificaMensagemCorreta(t, resultado, esperado)
    })


	t.Run("em espanhol", func(t *testing.T){
		resultado := Ola("Elodie","espanhol")
		esperado := "Hola, Elodie"
		verificaMensagemCorreta(t,resultado,esperado)
	})


	t.Run("em francês", func(t *testing.T){
		resultado := Ola("Mbape","frances")
		esperado := "Bonjour, Mbape"
		verificaMensagemCorreta(t,resultado,esperado)
	})
}

