package main

import (
	"os"

	"github.com/ajtroup1/goclear/src/lexer"
)

func main() {
	// Read in all the bytes from a hardcoded source file
	// TODO: Take in a param that denotes the file to run
		// ex. "make run repo/script.clr"
	bytes, _ := os.ReadFile("./examples/01.clr")

	// Extract all tokens (in a string) from the stringified byte array
	tokens := lexer.Tokenize(string(bytes))

	// Test print, displays every token for debugging
	for _, token := range tokens {
		token.Debug()
	}
}