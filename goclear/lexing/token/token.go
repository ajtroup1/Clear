/*

	Tokens are the structured building blocks of information that the Parser uses to generate the AST
	This file is used to define and reference the lowest level of source code information that can be
	referenced for building the AST and generating parsing errors

*/

package token

import (
	"fmt"
	"strconv"
)

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
	VOID    = "VOID"

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
	DOT      = "DOT"

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
	CONST    = "CONST"
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
	MODULE   = "MODULE"
	EXPORT   = "EXPORT"
	CONTINUE = "CONTINUE"
	BREAK    = "BREAK"
	WHILE    = "WHILE"
	FOR      = "FOR"
)

func (t *Token) Int() (int, error) {
	return strconv.Atoi(t.Literal)
}

func (t *Token) Float() (float64, error) {
	return strconv.ParseFloat(t.Literal, 64)
}

func (t *Token) Bool() (bool, error) {
	return strconv.ParseBool(t.Literal)
}

func (t *Token) Char() (rune, error) {
	if len(t.Literal) == 1 {
		return rune(t.Literal[0]), nil
	}

	unquoted, err := strconv.Unquote(t.Literal)
	if err != nil {
		return 0, fmt.Errorf("could not unquote character literal %q: %v", t.Literal, err)
	}

	if len(unquoted) != 1 {
		return 0, fmt.Errorf("character literal %q is not a single character", t.Literal)
	}

	return rune(unquoted[0]), nil
}
