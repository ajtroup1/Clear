package parser

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/goclear/lexing/token"
	"github.com/ajtroup1/goclear/parsing/ast"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.CONST:
		return p.parseConstStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	default:
		if isTypeToken(p.curToken.Type) {
			stmt := p.parseVariableDeclaration()
			if stmt.Type != token.ILLEGAL {
				return stmt
			}
		}
		return p.parseExpressionStatement()
	}
}

func isTypeToken(t token.TokenType) bool {
	typeTokens := []token.TokenType{
		token.INT,
		token.FLOAT,
		token.STRING,
		token.BOOL,
		token.CHAR,
	}
	for _, tt := range typeTokens {
		if t == tt {
			return true
		}
	}
	return false
}

func (p *Parser) parseVariableDeclaration() *ast.AssignStatement {
	stmt := &ast.AssignStatement{
		BaseNode: ast.BaseNode{Token: p.curToken},
		Type:     p.curToken.Type,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		BaseNode: ast.BaseNode{Token: p.curToken},
		Value:    p.curToken.Literal,
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		p.symbols[stmt.Name.Value] = strings.ToLower(string(stmt.Type)) 
		return stmt
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	fmt.Printf("stmt.Value: %v\n", stmt.Value.ToString())

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	rhsType := p.getTypeOfExpression(stmt.Value)
	if stmt.Type == token.ILLEGAL { 
		stmt.Type = token.TokenType(rhsType) 
		p.symbols[stmt.Name.Value] = strings.ToLower(rhsType)
	} else {
		if !p.checkTypeCompatibility(stmt.Value, string(stmt.Type)) {
			msg := fmt.Sprintf(
				"type mismatch: expected %s but got %s",
				strings.ToLower(string(stmt.Type)),
				rhsType,
			)
			p.errors = append(p.errors, ParserError{
				Msg: msg,
				Line: p.curToken.Line,
				Col: p.curToken.Col,
			})
			return &ast.AssignStatement{Type: token.ILLEGAL}
		}
	}

	p.symbols[stmt.Name.Value] = strings.ToLower(string(stmt.Type))

	return stmt
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	p.nextToken()

	myVar := p.parseVariableDeclaration()

	stmt.Name = myVar.Name
	stmt.Value = myVar.Value
	stmt.Type = myVar.Type

	return stmt
}

func (p *Parser) checkTypeCompatibility(expr ast.Expression, expectedType string) bool {
	actualType := p.getTypeOfExpression(expr)
	return actualType == strings.ToLower(expectedType)
}

func (p *Parser) getTypeOfExpression(expr ast.Expression) string {
	switch exp := expr.(type) {
	case *ast.Identifier:
		return p.symbols[exp.Value]
	case *ast.IntegerLiteral:
		return "int"
	case *ast.FloatLiteral:
		return "float"
	case *ast.BooleanLiteral:
		return "bool"
	case *ast.StringLiteral:
		return "string"
	case *ast.CharLiteral:
		return "char"
	case *ast.InfixExpression:
		leftType := p.getTypeOfExpression(exp.Left)
		rightType := p.getTypeOfExpression(exp.Right)
		if leftType == rightType {
			return leftType
		}
		return "type mismatch"
	case *ast.CallExpression:
		fmt.Printf("exp: %v\n", exp.ToString())
		fmt.Printf("**exp.FunctionIdentifier: %v\n", exp.FunctionIdentifier.ToString())
		fmt.Printf("symbols: %v\n", p.symbols)
		funcName := extractFunctionName(exp.FunctionIdentifier.ToString())
		fmt.Printf("funcName: %s\n", funcName)
		if retType, ok := p.symbols[funcName]; ok {
			return retType
		}
		return "unknown"
	default:
		return "unknown"
	}
}

func extractFunctionName(identifier string) string {
	identifier = strings.TrimSpace(identifier)
	identifier = strings.TrimPrefix(identifier, "IDENT ")
	identifier = strings.TrimSuffix(identifier, ":")
	return identifier
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	stmt.Expression = p.parseExpression(LOWEST)
	// fmt.Printf("stmt.Expression: %v\n", stmt.Expression.ToString())

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	if p.curTokenIs(token.SEMICOLON) {
		stmt.Init = nil
	} else if p.curTokenIs(token.LET) {
		stmt.Init = p.parseVariableDeclaration()
	} else {
		stmt.Init = p.parseExpressionStatement()
	}

	if !p.curTokenIs(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)
	if stmt.Condition == nil {
		return nil
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	stmt.Post = p.parsePostfixExpression(&ast.Identifier{BaseNode: ast.BaseNode{Token: p.peekToken}, Value: p.peekToken.Literal})
	if stmt.Post == nil {
		return nil
	}
	p.nextToken()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)
	if stmt.Condition == nil {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseModuleStatement() *ast.ModuleStatement {
	stmt := &ast.ModuleStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	// Expecting an identifier after 'import'
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken() // Consume the semicolon
	}

	return stmt
}
