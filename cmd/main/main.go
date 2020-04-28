package main

import (
	"fmt"

	pj "github.com/matsuyoshi30/go-pj"
)

func main() {
	sample := `{ "item": 42 }`

	tokenizer := pj.NewTokenizer(sample)
	fmt.Println("sample json is", sample)
	for _, t := range tokenizer.Tokenize() {
		fmt.Println("Name:", t.Name, ", Length:", t.Length)
	}
}
