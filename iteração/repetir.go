package iteracao

const quantidadeRepeticoes = 5
func Repetir(caractere string, quantidadeDeItertacoes int) string{
	var repeticoes string
	repetir := quantidadeRepeticoes

	if quantidadeDeItertacoes != 0 {
		repetir = quantidadeDeItertacoes
	}

	for i := 0; i < repetir; i++ {
		repeticoes += caractere
	}
	return repeticoes
}