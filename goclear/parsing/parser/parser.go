package parser

import (
	"fmt"

	"github.com/ajtroup1/goclear/lexing/lexer"
	"github.com/ajtroup1/goclear/lexing/token"
	"github.com/ajtroup1/goclear/parsing/ast"
	// "github.com/ajtroup1/goclear/utils"
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
	INDEX       // array[index]
	MEMBER      // object.member
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
	token.LBRACKET: INDEX,
	token.DOT:      MEMBER,
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

type Symbol struct {
	Name  string
	Type  ast.DataType
	Value interface{}
}

type Parser struct {
	l       *lexer.Lexer
	errors  ErrorList
	symbols map[string]Symbol

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

type ErrorList []ParserError
type ParserError struct {
	Msg  string
	Line int
	Col  int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:               l,
		errors:          ErrorList{},
		symbols:         make(map[string]Symbol),
		prefixParseFns:  make(map[token.TokenType]prefixParseFn),
		infixParseFns:   make(map[token.TokenType]infixParseFn),
		postfixParseFns: make(map[token.TokenType]postfixParseFn),
	}

	// Register prefix parse functions
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	// p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	// p.registerPrefix(token.CHAR, p.parseCharLiteral)
	// p.registerPrefix(token.TRUE, p.parseBoolean)
	// p.registerPrefix(token.FALSE, p.parseBoolean)
	// p.registerPrefix(token.BANG, p.parsePrefixExpression)
	// p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	// p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	// p.registerPrefix(token.IF, p.parseIfExpression)
	// p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// Register infix parse functions
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	// p.registerInfix(token.LPAREN, p.parseCallExpression)
	// p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	// p.registerInfix(token.DOT, p.parseMemberExpression)

	// Register postfix parse functions
	// p.registerPostfix(token.INC, p.parsePostfixExpression)
	// p.registerPostfix(token.DEC, p.parsePostfixExpression)

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

func (p *Parser) Errors() ErrorList {
	return p.errors
}

func (p *Parser) addError(msg string, line, col int) {
	p.errors = append(p.errors, ParserError{Msg: msg, Line: line, Col: col})
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
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

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, ParserError{Msg: msg, Line: p.peekToken.Line, Col: p.peekToken.Col})
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func mapTokenTypeToDataType(t token.TokenType) ast.DataType {
	switch t {
	case token.INT:
		return ast.INT
	case token.FLOAT:
		return ast.FLOAT
	case token.STRING:
		return ast.STRING
	case token.CHAR:
		return ast.CHAR
	case token.BOOL:
		return ast.BOOL
	case token.VOID:
		return ast.VOID
	default:
		return ast.UNKNOWN
	}
}

func (p *Parser) Parse() *ast.Program {
	return p.parseProgram()
}

func (p *Parser) parseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	// utils.PrettyPrintASTNode(program)

	return program
}
