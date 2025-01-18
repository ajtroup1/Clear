package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ajtroup1/clear/errors"
	"github.com/ajtroup1/clear/evaluator"
	"github.com/ajtroup1/clear/lexer"
	"github.com/ajtroup1/clear/logger"
	"github.com/ajtroup1/clear/modules"
	"github.com/ajtroup1/clear/object"
	"github.com/ajtroup1/clear/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	modules.Register(env)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		log := logger.NewLogger()

		line := scanner.Text()
		l := lexer.New(line, log, false)
		p := parser.New(l, log, false)

		program := p.ParseProgram()
		if len(p.Errors) != 0 {
			printParserErrors(out, p.Errors, program.String())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []errors.Error, line string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " line: "+line+"\n")
	io.WriteString(out, " parser errors:\n")
	for _, err := range errors {
		io.WriteString(out, "\t"+err.Message+"\n")
	}
}
