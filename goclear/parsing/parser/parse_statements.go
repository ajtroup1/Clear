package parser

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/goclear/lexing/token"
	"github.com/ajtroup1/goclear/parsing/ast"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.INT, token.FLOAT, token.STRING, token.CHAR, token.BOOL:
		return p.parseAssignStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{BaseNode: ast.BaseNode{Token: p.curToken}}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{BaseNode: ast.BaseNode{Token: p.curToken}}
	stmt.Type = mapTokenTypeToDataType(p.curToken.Type)

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
	if _, found := p.symbols[stmt.Name.Value]; found {
		p.addError(fmt.Sprintf("variable '%s' already declared", stmt.Name.Value), p.curToken.Line, p.curToken.Col)
	}

	if p.peekToken.Type == token.ASSIGN {
		p.nextToken()
		p.nextToken()

		stmt.Value = p.parseExpression(LOWEST)
		fmt.Printf("stmt.Value: %s\n", stmt.Value.GetType())

		if stmt.Type != stmt.Value.GetType() {
			p.addError(fmt.Sprintf(
				"type mismatch: cannot assign %s to %s", stmt.Value.GetType(), stmt.Type),
				p.curToken.Line, p.curToken.Col)
		}
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	p.symbols[stmt.Name.Value] = Symbol{
		Name:  stmt.Name.Value,
		Type:  stmt.Type,
		Value: stmt.Value,
	}

	return stmt
}

func (p *Parser) parseReassignStatement() *ast.ReassignStatement {
	stmt := &ast.ReassignStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	stmt.Name = &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
	if _, found := p.symbols[stmt.Name.Value]; !found {
		p.addError(fmt.Sprintf("variable '%s' not declared", stmt.Name.Value), p.curToken.Line, p.curToken.Col)
	}

	p.nextToken() // Consume '='
	p.nextToken() // Move to the expression

	stmt.Value = p.parseExpression(LOWEST)
	stringType := strings.Split(stmt.Value.ToString(), " ")[0]
	if p.symbols[stmt.Name.Value].Type != mapTokenTypeToDataType(token.TokenType(stringType)) {
		p.addError(fmt.Sprintf(
			"type mismatch: cannot assign %s to %s for var '%s'", stmt.Value.GetType(), p.symbols[stmt.Name.Value].Type, stmt.Name.Value),
			p.curToken.Line, p.curToken.Col)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	p.symbols[stmt.Name.Value] = Symbol{
		Name:  stmt.Name.Value,
		Type:  p.symbols[stmt.Name.Value].Type,
		Value: stmt.Value,
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	if p.peekToken.Type == token.ASSIGN {
		return p.parseReassignStatement()
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
