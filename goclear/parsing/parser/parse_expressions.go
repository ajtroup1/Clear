package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ajtroup1/goclear/lexing/token"
	"github.com/ajtroup1/goclear/parsing/ast"
	// "github.com/ajtroup1/goclear/utils"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		postfix := p.postfixParseFns[p.curToken.Type]
		if postfix == nil {
			p.noPrefixParseFnError(p.curToken.Type)
			return nil
		}

		return postfix(nil)
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		if infix, ok := p.infixParseFns[p.peekToken.Type]; ok {
			p.nextToken()
			leftExp = infix(leftExp)
		} else if postfix, ok := p.postfixParseFns[p.peekToken.Type]; ok {
			p.nextToken()
			leftExp = postfix(leftExp)
		} else {
			break
		}
	}

	// fmt.Printf("***Expression***: %s\n", leftExp.ToString())

	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
}

func (p *Parser) parseCharLiteral() ast.Expression {
	lit := p.curToken.Literal
	if len(lit) != 1 {
		msg := fmt.Sprintf("could not parse %q as char, chars must be one character only", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
		return nil
	}

	return &ast.CharLiteral{BaseNode: ast.BaseNode{Token: p.curToken}, Value: rune(lit[0])}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{BaseNode: ast.BaseNode{Token: p.curToken}}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{BaseNode: ast.BaseNode{Token: p.curToken}}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		BaseNode: ast.BaseNode{Token: p.curToken},
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parsePostfixExpression(lhs ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		BaseNode: ast.BaseNode{Token: p.curToken},
		Operator: p.curToken.Literal,
		Left:     lhs,
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		BaseNode: ast.BaseNode{Token: p.curToken},
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // Move past the LPAREN

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{BaseNode: ast.BaseNode{Token: p.curToken}}

	if p.peekToken.Type != token.LPAREN {
		msg := fmt.Sprintf("if condition must be closed within parentheses. Immediate char after 'if' must be '(', got '%s' instead", p.peekToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.peekToken.Line, Col: p.peekToken.Col})
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if p.curToken.Type != token.RPAREN {
		return nil
	}

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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{BaseNode: ast.BaseNode{Token: p.curToken}}

	if p.peekTokenIs(token.IDENT) {
		lit.Name = &ast.Identifier{BaseNode: ast.BaseNode{Token: p.peekToken}, Value: p.peekToken.Literal}
		p.nextToken()
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()
	for _, param := range lit.Parameters {
		p.symbols[param.Value] = strings.ToLower(string(param.Type))
	}

	p.nextToken()

	if (!p.curTokenIs(token.ASSIGN) || !p.curTokenIs(token.MINUS)) && !p.peekTokenIs(token.GT) {
		msg := fmt.Sprintf("expected '=>' or '->' after function parameters, got %s instead", p.peekToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.peekToken.Line, Col: p.peekToken.Col})
		return nil
	}

	p.nextToken()
	if !isTypeToken(p.peekToken.Type) {
		msg := fmt.Sprintf("expected type for function return after '=>', but got %s", p.peekToken.Literal)
		p.errors = append(p.errors, ParserError{Msg: msg, Line: p.peekToken.Line, Col: p.peekToken.Col})
		return nil
	}

	lit.ReturnType = p.peekToken.Type
	p.symbols[lit.Name.Value] = strings.ToLower(string(lit.ReturnType))

	p.nextToken()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// Handle empty parameter list
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// There is at least one parameter
	for {
		ident := &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}}

		if !p.curTokenIs(token.IDENT) {
			msg := fmt.Sprintf("expected identifier for parameter but got %s", p.curToken.Literal)
			p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
			return nil
		}
		ident.Value = p.curToken.Literal

		// Parse the parameter's type
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken() // Move past the COLON
		if !isTypeToken(p.curToken.Type) {
			msg := fmt.Sprintf("expected type for parameter but got %s", p.curToken.Literal)
			p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
			return nil
		}

		ident.Type = p.curToken.Type

		identifiers = append(identifiers, ident)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
			p.nextToken()
		} else {
			break
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{BaseNode: ast.BaseNode{Token: p.curToken}, FunctionIdentifier: function}
	exp.Arguments = p.parseCallArguments()

	// utils.PrettyPrintASTNode(exp)

	return exp
}

func (p *Parser) parseCallArguments() []ast.CallArgument {
	args := []ast.CallArgument{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()

	// Parse the first argument
	arg := &ast.CallArgument{BaseNode: ast.BaseNode{Token: p.curToken}}
	arg.Expression = p.parseExpression(LOWEST)

	_type := p.getTypeOfExpression(arg.Expression)
	arg.Type = stringToTokenType(_type)
	// fmt.Printf("Type: %s\n", arg.Type)

	args = append(args, *arg)

	// Parse any additional arguments
	for p.peekTokenIs(token.COMMA) {
		arg := &ast.CallArgument{BaseNode: ast.BaseNode{Token: p.curToken}}
		p.nextToken()
		p.nextToken()

		arg.Expression = p.parseExpression(LOWEST)

		_type := p.getTypeOfExpression(arg.Expression)
		arg.Type = stringToTokenType(_type)
		// fmt.Printf("Type: %s\n", arg.Type)

		args = append(args, *arg)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}