package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/ajtroup1/clear/src/evaluator"
	"github.com/ajtroup1/clear/src/lexer"
	"github.com/ajtroup1/clear/src/parser"
	"github.com/ajtroup1/clear/src/repl"
)

func main() {
	debug := flag.Bool("d", false, "Enable debug mode")
	flag.Parse()

	if flag.NArg() > 0 {
		filePath := flag.Arg(0)
		fmt.Printf("Running file: %s\n", filePath)
		runFile(filePath, *debug)
	} else {
		startREPL()
	}
}

func startREPL() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

func runFile(filePath string, debug bool) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Running script: %s\n", filePath)
	if debug {
		fmt.Println("Debug mode enabled")
	}
	// Example: pass content to the lexer and parser
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	if p.Errors() != "" {
		fmt.Println("Parser errors:")
		fmt.Printf("  %v\n", p.Errors())
		return
	}

	eval := evaluator.Eval(program)
	if eval != nil {
		fmt.Printf("Evaluation result: %v\n", eval.Inspect())
		return
	}

	if debug {
		fmt.Printf("Parsed program: %v\n", program.String())

	} else {
		fmt.Printf("Program parsed successfully...\n")
		return
	}

	astJSON, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling AST to JSON: %v\n", err)
		return
	}

	jsonFilePath := filePath + ".ast.json"
	err = os.WriteFile(jsonFilePath, astJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing AST JSON to file: %v\n", err)
		return
	}

	fmt.Printf("AST JSON written to %s\n", jsonFilePath)
}
