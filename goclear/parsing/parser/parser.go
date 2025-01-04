/*
	The parser is responsible for taking a stream of tokens and converting them into an abstract syntax tree (AST).
	The greatest tool the parser utilizes is recursion. The parser is able to recursively parse expressions, statements, and other constructs by calling itself.
	When the parser encounters a token, it will look up the appropriate parsing function to call based on the token type.
	- These parsing functions are handled through maps of token types to parsing functions.
	In the end, the parser will return a Program node that contains a slice of all the statements in the program.
*/

package parser

import (
	"fmt"

	"github.com/ajtroup1/goclear/lexing/lexer"
	"github.com/ajtroup1/goclear/lexing/token"
	"github.com/ajtroup1/goclear/parsing/ast"
)

// Tools to handle precedence levels and their association with token types
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

// Type aliases for parsing functions
type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

// Symbols are used to store variable names, types, and values
// It is necessary to keep track of variables via symbols:
// - To ensure that variables are not redeclared
// - To store the type of the variable for comparison during value reassignment
// - To store the value of the variable for reassignment, evaluation, etc.
type Symbol struct {
	Name  string
	Type  ast.DataType
	Value interface{}
}

// The overall structure of the parser
// Tracks information about the lexed tokens and parsing state
type Parser struct {
	l       *lexer.Lexer      // The lexer to pull tokens from
	errors  ErrorList         // List of Parser-specific errors
	symbols map[string]Symbol // Map of variable names to their Symbol objects

	// Token tracking for comparison and parsing information
	curToken  token.Token
	peekToken token.Token

	// Maps for parsing functions based on token type and number / order of expression elements
	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

// Error handling for the parser
type ErrorList []ParserError
type ParserError struct {
	Msg  string
	Line int
	Col  int
}

// Return the list of errors
// Since p.errors is unexported, a public method is necessary to access the errors
func (p *Parser) Errors() ErrorList {
	return p.errors
}

// Append an error to the list of errors
func (p *Parser) addError(msg string, line, col int) {
	p.errors = append(p.errors, ParserError{Msg: msg, Line: line, Col: col})
}

// Error handling for missing prefix parse functions
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, ParserError{Msg: msg, Line: p.curToken.Line, Col: p.curToken.Col})
}

// Create a new parser with a given lexer instance
func New(l *lexer.Lexer) *Parser {
	// Create a 'blank' parser with the given lexer
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
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.CHAR, p.parseCharLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// Register infix parse functions
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	// p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	// p.registerInfix(token.DOT, p.parseMemberExpression)

	// Register postfix parse functions
	p.registerPostfix(token.INC, p.parsePostfixExpression)
	p.registerPostfix(token.DEC, p.parsePostfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

// Helper functions for registering parsing functions to their respective token types
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

// Update the current and peek tokens
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Compare the current or peek token to a given token type
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Ensure that the next token is of a given type
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Error handling for unexpected tokens
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, ParserError{Msg: msg, Line: p.peekToken.Line, Col: p.peekToken.Col})
}

// Precedence handling for parsing infix expressions
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

// Return a corresponding AST data type for a given token type
func mapTokenTypeToDataType(t token.TokenType) ast.DataType {
	// fmt.Printf("t: %s\n", t)
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
	case token.IDENT:
		return ast.UNKNOWN
	default:
		return ast.UNKNOWN
	}
}

// Returns a bool indicating whether the current token is an AST data type
func (p *Parser) isDataType() bool {
	switch p.peekToken.Type {
	case token.INT, token.FLOAT, token.STRING, token.CHAR, token.BOOL:
		p.nextToken()
		return true
	default:
		return false
	}
}

// Public method to parse the program
func (p *Parser) Parse() *ast.Program {
	return p.parseProgram()
}

// Highest-level parsing function for the program
// Parses the program as a slice of statements and returns the program node
func (p *Parser) parseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}
    program.Imports = []*ast.ModuleStatement{}

    for p.curToken.Type != token.EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            if moduleStmt, ok := stmt.(*ast.ModuleStatement); ok {
                program.Imports = append(program.Imports, moduleStmt)
            } else {
                program.Statements = append(program.Statements, stmt)
            }
        }
        p.nextToken()
    }

    // utils.PrettyPrintASTNode(program)

    return program
}
