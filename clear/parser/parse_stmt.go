package parser

import (
	"github.com/ajtroup1/clear/ast"
	"github.com/ajtroup1/clear/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.MOD:
		return p.parseModuleStatement()
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionOrAssignStatement()
	}
}

func (p *Parser) parseModuleStatement() *ast.ModuleStatement {
	stmt := &ast.ModuleStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.COLON) {
		return nil
	}

	if p.peekTokenIs(token.ASTERISK) {
		stmt.ImportAll = true
		return stmt
	}

	if !p.expectPeek(token.LBRACKET) {
		return nil
	}

	if p.peekTokenIs(token.RBRACKET) {
		// TODO: Warning for empty import list
		return stmt
	}

	for !p.peekTokenIs(token.RBRACKET) {
		if !p.expectPeek(token.IDENT) {
			return nil
		}

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		stmt.Imports = append(stmt.Imports, ident)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionOrAssignStatement() ast.Statement {
	ident := p.parseIdentifier()
	if p.peekTokenIs(token.ASSIGN) {
		if id, ok := ident.(*ast.Identifier); ok {
			return p.parseAssignStatement(id)
		}
	}

	// fmt.Printf("current token: %s\n", p.curToken.Literal)
	return p.parseExpressionStatement()
}

func (p *Parser) parseAssignStatement(ident *ast.Identifier) *ast.AssignStatement {
	stmt := &ast.AssignStatement{Token: p.curToken, Name: ident}

	p.nextToken()
	p.nextToken()
	// fmt.Printf("current token: %s\n", p.curToken.Literal)

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
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
