func TestCorredor(t *testing.T) {
t.Run("compara a velocidade de servidores, retornando o endereço do mais rápido", func(t *testing.T) {
servidorLento := criarServidorComAtraso(20 _ time.Millisecond)
servidorRapido := criarServidorComAtraso(0 _ time.Millisecond)

        defer servidorLento.Close()
        defer servidorRapido.Close()

        URLLenta := servidorLento.URL
        URLRapida := servidorRapido.URL

        esperado := URLRapida
        resultado, err := Corredor(URLLenta, URLRapida)

        if err != nil {
            t.Fatalf("não esperava um erro, mas obteve um %v", err)
        }

        if resultado != esperado {
            t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
        }
    })

    t.Run("retorna um erro se o servidor não responder dentro de 10s", func(t *testing.T) {
        servidor := criarServidorComAtraso(25 * time.Millisecond)

        defer servidor.Close()

        _, err := Configuravel(servidor.URL, servidor.URL, 20*time.Millisecond)

        if err == nil {
            t.Error("esperava um erro, mas não obtive um")
        }
    })

}
