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
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	fmt.Printf("current token: %s\n", p.curToken.Type)

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

func (p *Parser) parseIdentifier() ast.Expression {
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