package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ajtroup1/clear/src/evaluator"
	"github.com/ajtroup1/clear/src/lexer"
	"github.com/ajtroup1/clear/src/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, err string) {
	io.WriteString(out, " parser errors:\n")
	io.WriteString(out, "" + err)
}