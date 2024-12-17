package parser

import (
	"github.com/ajtroup1/goclear/src/ast"
	"github.com/ajtroup1/goclear/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos int
}

func createParser (tokens []lexer.Token) *parser {
	createTokenLookups()
	// createTypeTokenLookups()

	p := &parser{
		tokens: tokens,
		pos: 0,
	}

	return p
}

func Parse (tokens []lexer.Token) ast.BlockStmt {
	p := createParser(tokens)
	body := make([]ast.Statement, 0)

	for p.hasTokens() {
		body = append(body, parse_stmt(p))
	}

	return ast.BlockStmt{
		Body: body,
	}
}