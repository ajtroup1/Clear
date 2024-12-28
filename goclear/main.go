package main

import (
	"fmt"
	"os"

	"github.com/ajtroup1/goclear/lexer"
)

func main() {
	debug := true
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path")
		return
	}

	filePath := os.Args[1]
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) > 2 && os.Args[2] == "-d" {
		debug = true
	}

	src := string(bytes)
	lexer := lexer.New(src)
	tokens := lexer.Lex()
	if debug {
		for i, token := range tokens {
			fmt.Println(i, token.Stringify())
		}
	}

}
