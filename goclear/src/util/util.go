package util

import (
	"fmt"

	"github.com/ajtroup1/goclear/src/types"
)

func PrintParserError(e types.ParserError) {
	fmt.Printf("Parser::Error --> (line %d, col %d): %s\n", e.Line, e.Column, e.Message)
	
	fmt.Printf("\t-- Line: %s\n", e.LineContent)
}
