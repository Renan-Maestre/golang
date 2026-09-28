package injecao

import (
	"bytes"
	"testing"
)

func TestCumprimenta(t *testing.T) {
    buffer := bytes.Buffer{}
    Cumprimenta(&buffer, "Chris")

    resultado := buffer.String()
    esperado := "Olá, Chris"

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }
}