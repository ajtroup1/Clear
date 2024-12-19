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
func printParserErrors(out io.Writer, errors []types.ParserError) {
	for _, pe := range errors {
		io.WriteString(out, fmt.Sprintf("\03333m\t (line %d, col %d) "+pe.Message+"\n\0330m", pe.Line, pe.Column))
	}
}
