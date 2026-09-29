package corredor

import (
	"fmt"
	"net/http"
	"time"
)


var limiteDeDezSegundos = 10 * time.Second

func Corredor(a, b string) (vencedor string, erro error) {
    return Configuravel(a, b, limiteDeDezSegundos)
}

//utilizando goroutines
func Configuravel(a,b string, tempoLimite time.Duration) (vencedor string, erro error){
	select{
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b,nil
	case <-time.After(tempoLimite):
		return "", fmt.Errorf("tempo limite de espera excedido para %s e %s", a, b)
	}
}

func ping(URL string) chan bool{
	ch := make(chan bool)

	go func() {
		http.Get(URL)
		ch <- true
	}()
	return ch
}


// implmentação normal
func Corredor2(a,b string)(vencedor string){
	duracaoA := medirTempoDeResposta(a)
	duracaoB :=  medirTempoDeResposta(b)

	if duracaoA < duracaoB {
		return a
	}

	return b
}


func medirTempoDeResposta(url string) time.Duration{
	inicio := time.Now()
	http.Get(url)
	return time.Since(inicio)
}