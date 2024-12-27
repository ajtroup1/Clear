package main

import (
	"fmt"
	"os"

	"github.com/ajtroup1/goclear/src/lexer"
	"github.com/ajtroup1/goclear/src/parser"
	"github.com/ajtroup1/goclear/src/repl"
)

func main() {
	replMode := true
	// printJson := true
	// dumpLitter := false

	if replMode {
		repl.Start(os.Stdin, os.Stdout)
	} else {
		filePath := "./examples/simpleScript.clr"
		srcBytes, err := os.ReadFile(filePath)
		if err != nil {
			panic(fmt.Sprintf("File '%s' not found", filePath))
		}
		fmt.Printf("Reading in file '%s'\n", filePath)
		src := string(srcBytes)

		lexer := lexer.New(src)
		parser := parser.New(lexer)
		program := parser.ParseProgram()

		fmt.Printf("%s", program.String())
	}
}