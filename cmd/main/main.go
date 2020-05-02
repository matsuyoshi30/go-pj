package main

import (
	"fmt"
	"os"

	pj "github.com/matsuyoshi30/go-pj"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <json>\n", os.Args[0])
		return
	}

	input := os.Args[1]
	l := pj.NewLexer(input)
	p := pj.NewParser(l)

	root, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("=> ")
	root.PrintFromRoot()
}
