package util

import (
	"fmt"

	"github.com/ajtroup1/goclear/src/types"
)

func PrintParserError(e types.ParserError) {
	fmt.Printf("\033[31mParser::Error --> (line %d, col %d): %s\n\033[0m", e.Line, e.Column, e.Message)

	fmt.Printf("\t\033[33m-- Line: %s\n\033[0m", e.LineContent)
}
