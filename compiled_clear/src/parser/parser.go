package parser

import (
	"fmt"

	"github.com/ajtroup1/compiled_clear/src/ast"
	"github.com/ajtroup1/compiled_clear/src/errorlogger"
	"github.com/ajtroup1/compiled_clear/src/lexer"
	"github.com/ajtroup1/compiled_clear/src/token"
)

const (
	_ int = iota
	LOWEST
	OR          // ||
	AND         // &&
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	POSTFIX     // x++ or x--
	CALL        // myFunction(x)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.PLUS_EQ:  SUM,
	token.MINUS_EQ: SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MULT_EQ:  PRODUCT,
	token.DIV_EQ:   PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.INC:      POSTFIX,
	token.DEC:      POSTFIX,
	token.AND:      AND,
	token.OR:       OR,
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	el        *errorlogger.ErrorLogger
	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

func New(l *lexer.Lexer, el *errorlogger.ErrorLogger) *Parser {
	p := &Parser{
		l:               l,
		el:              el,
		prefixParseFns:  make(map[token.TokenType]prefixParseFn),
		infixParseFns:   make(map[token.TokenType]infixParseFn),
		postfixParseFns: make(map[token.TokenType]postfixParseFn),
	}

	// Prefix parse function handlers
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	// Infix parse function handlers
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.PLUS_EQ, p.parseInfixExpression)
	p.registerInfix(token.MINUS_EQ, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MULT_EQ, p.parseInfixExpression)
	p.registerInfix(token.DIV_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	// p.registerInfix(token.LPAREN, p.parseCallExpression)
	// p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// Postfix parse function handlers
	// p.registerInfix(token.INC, p.parsePostfixExpression)
	// p.registerInfix(token.DEC, p.parsePostfixExpression)

	p.nextToken()
	p.nextToken()

	return p
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

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
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

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.el.NewError(p.peekToken.Line, p.curToken.Col, fmt.Sprintf("expected next token to be %s, got %s ('%s') instead", t, p.peekToken.Type, p.peekToken.Literal), "Parsing", false)
	return false
}
