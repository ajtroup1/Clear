package parser

import (
	"fmt"

	"github.com/ajtroup1/goclear/lexing/token"
	"github.com/ajtroup1/goclear/parsing/ast"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

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

func (p *Parser) parseStringLiteral() ast.Expression {
		return &ast.StringLiteral{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}