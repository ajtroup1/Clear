package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ajtroup1/goclear/lexer"
	"github.com/ajtroup1/goclear/parser"
	"github.com/sanity-io/litter"
)

func main() {
	jsonMode := true
	litterMode := false
	debug := false
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
	lexer.Lex()
	if debug {
		for i, token := range lexer.Tokens {
			fmt.Println(i, token.Stringify())
		}
	}

	parser := parser.New(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) != 0 {
		fmt.Println("Parser errors:")
		for _, err := range parser.Errors() {
			fmt.Printf("\tParser::Error --> '%s' [line: %d, col: %d]\n", err.Msg, err.Line, err.Col)
		}
		return
	}

	if debug {
		if len(program.Statements) == 0 {
			fmt.Println("No program statements")
			return
		}
		fmt.Println("Program Statements:")
		if jsonMode {
			file, err := os.Create("program.json")
			if err != nil {
				fmt.Println("Error creating JSON file:", err)
				return
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(program); err != nil {
				fmt.Println("Error encoding program to JSON:", err)
				return
			}
			fmt.Println("Program written to program.json")
		}
		if litterMode {
			litter.Dump(program)
		}
	}
}
