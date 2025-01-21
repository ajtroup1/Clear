package token

import (
	"fmt"

	"github.com/ajtroup1/clear/logger"
)

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"   // 1343456
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	INC      = "++"
	DEC      = "--"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	DOT       = "."
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	MOD      = "MOD"
	WHILE    = "WHILE"
	FOR      = "FOR"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

func (t *Token) String() string {
	return fmt.Sprintf("- Token::%s '%s' at [line: %d, col: %d]", t.Type, t.Literal, t.Line, t.Col)
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"mod":    MOD,
	"module": MOD,
	"while":  WHILE,
	"for":    FOR,
}

func LookupIdent(ident string, logger *logger.Logger, enc int) TokenType {
	if tok, ok := keywords[ident]; ok {
		logger.Append(fmt.Sprintf("%d. Discerned that '%s' is a keyword '%s'\n", enc, ident, tok))
		return tok
	}
	return IDENT
}
