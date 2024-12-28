package token

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

func (t *Token) Stringify() string {
	return fmt.Sprintf("%s ('%s') [line: %d, col: %d]", t.Type, t.Literal, t.Line, t.Col)
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	NULL    = "NULL"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"   // 1343456
	FLOAT  = "FLOAT" // 134.3456
	STRING = "STRING"
	CHAR   = "CHAR"
	BOOL   = "BOOL"

	// Operators
	ASSIGN   = "ASSIGN"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	BANG     = "BANG"
	ASTERISK = "ASTERISK"
	SLASH    = "SLASH"
	LT       = "LT"
	LT_EQ    = "LT_EQ"
	GT       = "GT"
	GT_EQ    = "GT_EQ"
	EQ       = "EQ"
	NOT_EQ   = "NOT_EQ"
	PLUS_EQ  = "PLUS_EQ"
	MINUS_EQ = "MINUS_EQ"
	MUL_EQ   = "MUL_EQ"
	DIV_EQ   = "DIV_EQ"
	INC      = "INC"
	DEC      = "DEC"

	// Logical Operators
	AND = "AND"
	OR  = "OR"

	// Delimiters
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"
	COLON     = "COLON"

	LPAREN   = "LPAREN"
	RPAREN   = "RPAREN"
	LBRACE   = "LBRACE"
	RBRACE   = "RBRACE"
	LBRACKET = "LBRACKET"
	RBRACKET = "RBRACKET"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	NEW      = "NEW"
	CLASS    = "CLASS"
	THIS     = "THIS"
	SUPER    = "SUPER"
	STATIC   = "STATIC"
	IMPORT   = "IMPORT"
	EXPORT   = "EXPORT"
	CONTINUE = "CONTINUE"
	BREAK    = "BREAK"
	WHILE    = "WHILE"
	FOR      = "FOR"
)
