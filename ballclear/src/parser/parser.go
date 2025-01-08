package parser

import (
	"bytes"
	"fmt"

	"github.com/ajtroup1/clear/src/ast"
	"github.com/ajtroup1/clear/src/lexer"
	"github.com/ajtroup1/clear/src/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.INC:      PREFIX,
	token.DEC:      PREFIX,
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []ParserError

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []ParserError{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.INC, p.parsePostfixExpression)
	p.registerPostfix(token.DEC, p.parsePostfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) expectCurrent(t token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.curError(t)
		return false
	}
}



type ParserError struct {
	Token token.Token
	Msg   string
}

func (p *Parser) Errors() string {
	var out bytes.Buffer

	if len(p.errors) == 0 {
		return ""
	}

	// Header with formatting and separators
	out.WriteString("\n-----------------------------------\n\n")
	out.WriteString(fmt.Sprintf("\033[31mParsing Errors (%d):\033[0m\n", len(p.errors)))

	// Loop through errors and format each one
	for _, err := range p.errors {
		out.WriteString(fmt.Sprintf(
			"\t\033[31mParser::Error ---> \"%s\" [line: %d, col: %d]\033[0m\n",
			err.Msg, err.Token.Line, err.Token.Col))
	}

	// Footer separator
	out.WriteString("\n-----------------------------------\n\n")

	return out.String()
}

func (p *Parser) peekError(t token.TokenType) {
	pe := ParserError{Token: p.peekToken, Msg: fmt.Sprintf("expected next token to be '%s', got '%s' instead", t, p.peekToken.Type)}
	p.errors = append(p.errors, pe)
}

func (p *Parser) curError(t token.TokenType) {
	pe := ParserError{Token: p.curToken, Msg: fmt.Sprintf("expected current token to be '%s', got '%s' instead", t, p.curToken.Type)}
	p.errors = append(p.errors, pe)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	pe := ParserError{Token: p.curToken, Msg: fmt.Sprintf("no prefix parse function for '%s' found", t)}
	p.errors = append(p.errors, pe)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}
