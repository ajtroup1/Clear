package parser

import (
	"fmt"

	"github.com/ajtroup1/clear/src/ast"
	"github.com/ajtroup1/clear/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.MODULE:
		return p.parseModuleStatement()
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
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
		p.nextToken() 
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	}

	if !p.expectPeek(token.LBRACKET) {
		return nil
	}

	if p.peekTokenIs(token.RBRACKET) {
		fmt.Println("Warning: Empty import list for module", stmt.Name.Value)
		p.nextToken() 
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken() 
		}
		return stmt
	}

	for !p.peekTokenIs(token.RBRACKET) {
		if p.peekTokenIs(token.EOF) {
			p.errors = append(p.errors, ParserError{Token: p.peekToken, Msg: "expected ']' to close import list"})
			return nil
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken() 
			continue
		}

		imp := &ast.Identifier{Token: p.peekToken, Value: p.peekToken.Literal}
		stmt.Imports = append(stmt.Imports, imp)
		p.nextToken()
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.peekTokenIs(token.ASSIGN) {
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	} else if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
	} else {
		// !
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

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		stmt.Body = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		stmt.Body = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	p.nextToken()

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
	}

	if !p.expectPeek(token.LET) {
		return nil
	}
	stmt.Init = p.parseLetStatement()

	
	if !p.expectCurrent(token.SEMICOLON) {
		return nil
	}
	
	stmt.Condition = p.parseExpression(LOWEST)
	
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	
	p.nextToken()
	stmt.Post = p.parseExpression(LOWEST)

	
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{Token: p.curToken}

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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
