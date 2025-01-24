package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ajtroup1/compiled_clear/src/errorlogger"
	"github.com/ajtroup1/compiled_clear/src/lexer"
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
		// startRepl()
	} else {
		fmt.Println("Error: Invalid arguments")
		os.Exit(1)
	}
}

func runScript(filePath string, debug bool) {
	if debug {
		fmt.Println("Running script in debug mode")
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	src := string(file)
	lines := strings.Split(src, "\n")

	el := errorlogger.New(lines, debug)

	l := lexer.New(src, el, debug)
	fmt.Printf("Lexer: %v\n", l)
	// p := parser.New(l)

	// program := p.ParseProgram()
	// if len(p.Errors()) != 0 {
	// 	fmt.Println("Error: Parsing failed")
	// 	for _, msg := range p.Errors() {
	// 		fmt.Printf("\t%s\n", msg)
	// 	}
	// 	os.Exit(1)
	// }

	// env := object.NewEnvironment()
	// eval := evaluator.New(program, env)

	// if debug {
	// 	fmt.Println("### Program:")
	// 	fmt.Println(program.String())
	// 	fmt.Println("### Environment:")
	// 	fmt.Println(env.String())
	// }

	// if err := eval.Run(); err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// 	os.Exit(1)
	// }
}
