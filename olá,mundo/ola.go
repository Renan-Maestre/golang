package main

import "fmt"

const prefixoOlaPortugues = "Olá, "


func Ola(nome string) string {
	if nome == ""{
		nome = "Mundo"
	}
	
	return prefixoOlaPortugues + nome
}


func main() {
	fmt.Println(Ola("Renan"))
}