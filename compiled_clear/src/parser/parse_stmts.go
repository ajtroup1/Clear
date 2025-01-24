package parser

import (
	"fmt"

	"github.com/ajtroup1/compiled_clear/src/ast"
	"github.com/ajtroup1/compiled_clear/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	fmt.Print() // Disables unused import error
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken} // token.LET

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()

	expr := p.parseExpression(LOWEST)
	stmt.Value = expr

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken} // token.RETURN

	p.nextToken()

	expr := p.parseExpression(LOWEST)
	stmt.ReturnValue = expr

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}