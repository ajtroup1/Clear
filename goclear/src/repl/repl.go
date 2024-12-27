package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ajtroup1/goclear/src/lexer"
	"github.com/ajtroup1/goclear/src/parser"
	"github.com/ajtroup1/goclear/src/types"
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
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

const MONKEY_FACE = ` __,__
  .--. .-"     "-. .--.
/ .. \/  .-. .-.  \/ .. \
| |  '|  /   Y   \  |'  |
| \   \  \ 0 | 0 /  /   /
 \ '- ,\.-"""""""-./, -' /
  ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []types.ParserError) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, pe := range errors {
		io.WriteString(out, "\t"+pe.Message+"\n")
	}
}
