package inteiros

import "fmt"

// Adiciona recebe dois inteiros e retorna a soma deles
func Adicionar(x,y int)int{
	return x+y
}

func main(){
	fmt.Println(Adicionar(2,2))
}