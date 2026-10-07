package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	argValue := os.Args[1]
	argType := os.Args[2]

	valor, err := strconv.ParseFloat(argValue, 64)
	if err != nil {
		fmt.Println("Não foi possivel converter o input para número : ", err)
	}

	fmt.Println(conversor(valor, argType))
}

func conversor(value float64, moeda string) float64 {
	rates := Rates
	return value * rates[moeda]
}
