package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/ajtroup1/clear/evaluator"
	"github.com/ajtroup1/clear/object"
	"github.com/ajtroup1/clear/parser"

	"github.com/ajtroup1/clear/lexer"
	"github.com/ajtroup1/clear/repl"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.BoolVar(&debug, "d", false, "Debug mode (short)")
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		filepath := args[0]

		if !strings.HasSuffix(filepath, ".clr") {
			fmt.Println("Error: Invalid file type. Please provide a .clr file")
			os.Exit(1)
		}

		runScript(filepath, debug)
	} else if len(args) == 0 {
		startRepl()
	} else {
		fmt.Println("Error: Invalid arguments")
		os.Exit(1)
	}
}

func startRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Clear programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

func runScript(filepath string, debug bool) {
	if debug {
		fmt.Printf("Executing \"%s\"\n", filepath)
	}

	// Read the source file
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	src := string(bytes)
	lexer := lexer.New(src)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	if debug {
		// Generate JSON representation of the parse tree
		parseTreeJSON, err := json.MarshalIndent(program, "", "  ")
		if err != nil {
			fmt.Printf("Error generating parse tree JSON: %s\n", err)
			os.Exit(1)
		}

		// Construct the output file path
		jsonFilePath := strings.TrimSuffix(filepath, ".clr") + ".ast.json"

		err = os.WriteFile(jsonFilePath, parseTreeJSON, 0644)
		if err != nil {
			fmt.Printf("Error writing parse tree JSON to file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Parse tree JSON dumped to: %s\n", jsonFilePath)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	fmt.Printf("Evaluated: %s\n", evaluated.Inspect())
}
