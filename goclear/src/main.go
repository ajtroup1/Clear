package main

import (
	"os"

	"github.com/ajtroup1/goclear/src/lexer"
)

func main() {
	bytes, _ := os.ReadFile("./examples/01.clr")

	tokens := lexer.Tokenize(string(bytes))

	for _, token := range tokens {
		token.Debug()
	}
}