package main

import "fmt"

const espanhol = "espanhol"
const frances = "frances"
const prefixoOlaPortugues = "Olá, "
const prefixoOlaEspanhol = "Hola, "
const prefixoOlaFrances = "Bonjour, "

func main() {
	fmt.Println(Ola("Renan", "frances"))
}


func Ola(nome string, idioma string) string {
	if nome == ""{
		nome = "Mundo"
	}
	
	
	return prefixoSaudacao(idioma) + nome
}


func prefixoSaudacao(idioma string) (prefixo string){
	switch idioma {
	case frances:
		prefixo = prefixoOlaFrances
	case espanhol:
		prefixo = prefixoOlaEspanhol
	default:
		prefixo = prefixoOlaPortugues
	}
	return
}