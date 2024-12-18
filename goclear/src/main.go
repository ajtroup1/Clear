package main

import (
	// "encoding/json"
	// "fmt"
	"os"

	"github.com/ajtroup1/goclear/src/lexer"
	"github.com/ajtroup1/goclear/src/parser"
	"github.com/sanity-io/litter"
)

func main() {
	srcBytes, _ := os.ReadFile("./examples/00.clr")
	src := string(srcBytes)

	lexer := lexer.New(src)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	litter.Dump(program)

	// programJSON, err := json.MarshalIndent(program, "", "  ")
	// if err != nil {
	// 	fmt.Printf("Error serializing program to JSON: %v\n", err)
	// 	return
	// }

	// err = os.WriteFile("program.json", programJSON, 0644)
	// if err != nil {
	// 	fmt.Printf("Error writing JSON to file: %v\n", err)
	// 	return
	// }

	// fmt.Println("Program has been dumped to program.json")
}
