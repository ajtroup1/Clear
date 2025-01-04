// Outsource the parsing of expressions to separate file.

package parser

import (
	"fmt"

	"github.com/ajtroup1/goclear/lexing/token"
	"github.com/ajtroup1/goclear/parsing/ast"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.curToken == (token.Token{}) {
		p.addError("unexpected end of file", 0, 0)
		return nil
	}

	// Handle prefix expressions (e.g., identifiers, literals, negations, etc.)
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	// Handle postfix expressions (e.g., x++, y--).
	// If the next token is a valid postfix operator, process it.
	for p.peekTokenIs(token.INC) || p.peekTokenIs(token.DEC) {
		postfix := p.postfixParseFns[p.peekToken.Type]
		if postfix != nil {
			p.nextToken()              // Move past the postfix operator (e.g., '++', '--')
			leftExp = postfix(leftExp) // Apply postfix operation
		}
	}

	// Handle infix expressions based on precedence
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{BaseNode: ast.BaseNode{Token: p.curToken}, Operator: p.curToken.Literal}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{BaseNode: ast.BaseNode{Token: p.curToken}, Operator: p.curToken.Literal, Left: left}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.PostfixExpression{BaseNode: ast.BaseNode{Token: p.curToken}, Operator: p.curToken.Literal, Left: left}

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	found, ok := p.symbols[p.curToken.Literal]
	if ok {
		return &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal, Type: found.Type}
	}
	return &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{BaseNode: ast.BaseNode{Token: p.curToken}}

	value, err := p.curToken.Int()
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
		return nil
	}

	lit.Value = int64(value)
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{BaseNode: ast.BaseNode{Token: p.curToken}}

	value, err := p.curToken.Float()
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}

func (p *Parser) parseCharLiteral() ast.Expression {
	value, err := p.curToken.Char()
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as char", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
		return nil
	}

	return &ast.CharLiteral{BaseNode: ast.BaseNode{Token: p.curToken}, Value: value}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	value, err := p.curToken.Bool()
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as bool", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
		return nil
	}

	return &ast.BooleanLiteral{BaseNode: ast.BaseNode{Token: p.curToken}, Value: value}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{BaseNode: ast.BaseNode{Token: p.curToken}}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	p.nextToken()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{BaseNode: ast.BaseNode{Token: p.curToken}}
	p.nextToken()
	lit.Name = p.parseIdentifier().(*ast.Identifier)

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.ARROW) {
		return nil
	}

	if !p.isDataType() {
		return nil
	}

	lit.ReturnType = mapTokenTypeToDataType(p.curToken.Type)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	p.checkReturnTypeConsistency(lit)

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}

	if !p.expectPeek(token.COLON) {
		return nil
	}

	if !p.isDataType() {
		return nil
	}

	ident.Type = mapTokenTypeToDataType(p.curToken.Type)

	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
		if !p.expectPeek(token.COLON) {
			return nil
		}

		if !p.isDataType() {
			return nil
		}

		ident.Type = mapTokenTypeToDataType(p.curToken.Type)

		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) checkReturnTypeConsistency(lit *ast.FunctionLiteral) {
	for _, stmt := range lit.Body.Statements {
		if retStmt, ok := stmt.(*ast.ReturnStatement); ok {
			if retStmt.Value.GetType() != "" {
				// Simply inferred type such as INT, FLOAT, STRING, etc.
				if retStmt.Value.GetType() != lit.ReturnType {
					p.addError(fmt.Sprintf(
						"funciton return type mismatch for '%s': expected %s, got %s",
						lit.Name.Value, lit.ReturnType, retStmt.Value.GetType()),
						retStmt.Token.Line, retStmt.Token.Col)
				}
			}
		}
	}
}
