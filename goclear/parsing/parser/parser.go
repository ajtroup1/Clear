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
	POSTFIX     // X++ or X--
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
}

type SymbolTable map[string]string

func NewSymbolTable() SymbolTable {
	return make(SymbolTable)
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

type ParserError struct {
	Msg  string
	Line int
	Col  int
}

type Parser struct {
	l      *lexer.Lexer
	AST    *ast.Program
	errors []ParserError

	curToken  token.Token
	peekToken token.Token
	symbols   SymbolTable

	pos int

	imports        []*ast.ModuleStatement
	parsingImports bool

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		parsingImports: false,
		errors:         []ParserError{},
		symbols:        NewSymbolTable(),
	}

	p.AST = &ast.Program{Statements: []ast.Statement{}, Imports: []*ast.ModuleStatement{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.CHAR, p.parseCharLiteral)

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
	p.curToken = l.Tokens[0]
	p.peekToken = l.Tokens[1]
	p.pos = 0

	return p
}

func (p *Parser) nextToken() {
	p.pos++
	if p.pos >= len(p.l.Tokens) {
		p.curToken = token.Token{Type: token.EOF, Literal: "", Line: 0, Col: 0}
		p.peekToken = token.Token{Type: token.EOF, Literal: "", Line: 0, Col: 0}
		return
	}

	p.curToken = p.l.Tokens[p.pos]

	if p.pos+1 < len(p.l.Tokens) {
		p.peekToken = p.l.Tokens[p.pos+1]
	} else {
		p.peekToken = token.Token{Type: token.EOF, Literal: "", Line: 0, Col: 0}
	}
}

func (p *Parser) regressToken() {
	p.pos--
	if p.pos < 0 {
		p.curToken = token.Token{Type: token.EOF, Literal: "", Line: 0, Col: 0}
		p.peekToken = token.Token{Type: token.EOF, Literal: "", Line: 0, Col: 0}
		return
	}

	p.curToken = p.l.Tokens[p.pos]

	if p.pos+1 < len(p.l.Tokens) {
		p.peekToken = p.l.Tokens[p.pos+1]
	} else {
		p.peekToken = token.Token{Type: token.EOF, Literal: "", Line: 0, Col: 0}
	}
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

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, ParserError{Msg: msg, Line: p.peekToken.Line, Col: p.peekToken.Col})
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
		Imports:    []*ast.ModuleStatement{},
	}

	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.MODULE) && !p.parsingImports {
			p.parsingImports = true
		}

		if p.parsingImports {
			if p.curTokenIs(token.MODULE) {
				stmt := p.parseModuleStatement()
				if stmt != nil {
					program.Imports = append(program.Imports, stmt)
					p.AST.Imports = append(p.AST.Imports, stmt)
				}
			} else {
				p.parsingImports = false
			}
		}

		if !p.parsingImports {
			stmt := p.parseStatement()
			if stmt != nil {
				program.Statements = append(program.Statements, stmt)
			}
		}

		p.nextToken()
	}

	return program
}

// Define the function to convert string type to token.TokenType
func stringToTokenType(typeStr string) token.TokenType {
	switch typeStr {
	case "int":
		return token.INT
	case "float":
		return token.FLOAT
	case "string":
		return token.STRING
	case "char":
		return token.CHAR
	case "bool":
		return token.BOOL
	default:
		return token.ILLEGAL
	}
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
