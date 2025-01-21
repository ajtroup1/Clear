package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ajtroup1/clear/errors"
	"github.com/ajtroup1/clear/evaluator"
	"github.com/ajtroup1/clear/logger"
	"github.com/ajtroup1/clear/modules"
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
		filePath := args[0]

		if !strings.HasSuffix(filePath, ".clr") {
			fmt.Println("Error: Invalid file type. Please provide a .clr file")
			os.Exit(1)
		}

		runScript(filePath, debug)
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

func runScript(filePath string, debug bool) {
	if debug {
		fmt.Printf("Executing \"%s\"\n", filePath)
	}

	// Read the source file
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger()

	fileName := filepath.Base(filePath)

	if debug {
		log.InitText(fileName)
	}

	src := string(bytes)
	lexer := lexer.New(src, log, debug)
	parser := parser.New(lexer, log, debug)
	program := parser.ParseProgram()
	if program.NoStatements {
		fmt.Println("No valid statements in the program")
		os.Exit(1)
	}

	if debug {
		// Generate JSON representation of the parse tree
		parseTreeJSON, err := json.MarshalIndent(program, "", "  ")
		if err != nil {
			fmt.Printf("Error generating parse tree JSON: %s\n", err)
			os.Exit(1)
		}

		// Construct the output file path
		jsonfilePath := strings.TrimSuffix(filePath, ".clr") + ".ast.json"

		err = os.WriteFile(jsonfilePath, parseTreeJSON, 0644)
		if err != nil {
			fmt.Printf("Error writing parse tree JSON to file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Parse tree JSON dumped to: %s\n", jsonfilePath)
	}

	env := object.NewEnvironment()
	modules.Register(env)
	evaluator.Init(log, debug, lexer.Lines)
	evaluated := evaluator.Eval(program, env)

	errs, warn := errors.HasErrors(lexer.Errors, parser.Errors)
	if errs {
		fmt.Print(errors.ReportErrors(lexer.Errors, parser.Errors))
		os.Exit(1)
	}
	if warn {
		fmt.Print(errors.ReportErrors(lexer.Errors, parser.Errors))
	}

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Print(errors.ReportEvaluationError(evaluated.(*object.Error)))
		os.Exit(1)
	}

	if evaluated == nil {
		fmt.Println("\nProgram returned 0 (default)\n")
		os.Exit(0)
	}

	fmt.Printf("\nProgram returned: %s\n\n", evaluated.Inspect())

	if debug {
		logfilePath := strings.TrimSuffix(filePath, ".clr") + ".log.md"

		log.WriteFile(logfilePath)
		fmt.Printf("Log dumped to: %s\n", logfilePath)
	}
}
