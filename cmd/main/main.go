package main

import (
	"fmt"

	pj "github.com/matsuyoshi30/go-pj"
)

func main() {
	sample := `{ "item": 42 }`
	fmt.Println("sample json is", sample)

	l := pj.NewLexer(sample)
	p := pj.NewParser(l)

	root, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("=> ")
	root.PrintFromRoot()
}
