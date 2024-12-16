package lexer

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	// Special Tokens
	EOF     TokenType = iota
	ILLEGAL           // Token for unrecognized input
	NULL              //

	// Identifiers & Literals
	IDENT   // x, foo, myVariable
	NUMBER  // Numeric literals
	STRING  // String literals
	BOOLEAN // True, False, etc.

	// Punctuation
	COMMA    // ,
	DOT      // .
	DOT_DOT  // ..
	COLON    // :
	SEMI     // ;
	QUESTION // ?

	// Delimiters
	OPEN_BRACE    // {
	CLOSE_BRACE   // }
	OPEN_PAREN    // (
	CLOSE_PAREN   // )
	OPEN_BRACKET  // [
	CLOSE_BRACKET // ]

	// Operators
	PLUS        // +
	PLUS_PLUS   // ++
	MINUS       // -
	MINUS_MINUS // --
	STAR        // *
	SLASH       // /
	PERCENT     // %

	// Assignment Operators
	ASSIGNMENT  // =
	PLUS_EQUAL  // +=
	MINUS_EQUAL // -=
	STAR_EQUAL  // *=
	SLASH_EQUAL // /=

	// Comparison Operators
	BANG          // !
	COMPARISON    // ==
	NOT_EQUAL     // !=
	LESS          // <
	GREATER       // >
	LESS_EQUAL    // <=
	GREATER_EQUAL // >=
	TRUE
	FALSE

	// Logical Operators
	AND // &&
	OR  // ||

	// Keywords
	LET
	CONST
	CLASS
	NEW
	IMPORT
	FROM
	FN
	IF
	ELSE
	FOR
	FOREACH
	WHILE
	EXPORT
	TYPEOF
	IN
)

var keyword_lookup map[string]TokenType = map[string]TokenType{
	"true":    TRUE,
	"false":   FALSE,
	"null":    NULL,
	"let":     LET,
	"const":   CONST,
	"class":   CLASS,
	"new":     NEW,
	"import":  IMPORT,
	"from":    FROM,
	"fn":      FN,
	"if":      IF,
	"else":    ELSE,
	"foreach": FOREACH,
	"while":   WHILE,
	"for":     FOR,
	"export":  EXPORT,
	"typeof":  TYPEOF,
	"in":      IN,
}

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) isOneOfMany(expectedTokens ...TokenType) bool {
	for _, expected := range expectedTokens {
		if expected == t.Type {
			return true
		}
	}

	return false
}

func (t Token) Debug() {
	if t.isOneOfMany(IDENT, NUMBER, STRING) {
		fmt.Printf("%s (%s)\n", strings.ToUpper(TokenTypeString(t.Type)), t.Literal)
	} else {
		fmt.Printf("%s ()\n", strings.ToUpper(TokenTypeString(t.Type)))
	}
}

func NewToken(t TokenType, lit string) Token {
	return Token{
		Type:    t,
		Literal: lit,
	}
}

func TokenTypeString(t TokenType) string {
	switch t {
	case EOF:
		return "eof"
	case NULL:
		return "null"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case IDENT:
		return "identifier"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_BRACE:
		return "open_curly"
	case CLOSE_BRACE:
		return "close_curly"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case ASSIGNMENT:
		return "assignment"
	case COMPARISON:
		return "equals"
	case NOT_EQUAL:
		return "not_equals"
	case BANG:
		return "not"
	case LESS:
		return "less"
	case LESS_EQUAL:
		return "less_equals"
	case GREATER:
		return "greater"
	case GREATER_EQUAL:
		return "greater_equals"
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "dot"
	case DOT_DOT:
		return "dot_dot"
	case SEMI:
		return "semi_colon"
	case COLON:
		return "colon"
	case QUESTION:
		return "question"
	case COMMA:
		return "comma"
	case PLUS_PLUS:
		return "plus_plus"
	case MINUS_MINUS:
		return "minus_minus"
	case PLUS_EQUAL:
		return "plus_equals"
	case MINUS_EQUAL:
		return "minus_equals"
	case PLUS:
		return "plus"
	case MINUS:
		return "minus"
	case SLASH:
		return "slash"
	case STAR:
		return "star"
	case PERCENT:
		return "percent"
	case LET:
		return "let"
	case CONST:
		return "const"
	case CLASS:
		return "class"
	case NEW:
		return "new"
	case IMPORT:
		return "import"
	case FROM:
		return "from"
	case FN:
		return "fn"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case FOREACH:
		return "foreach"
	case FOR:
		return "for"
	case WHILE:
		return "while"
	case EXPORT:
		return "export"
	case IN:
		return "in"
	default:
		return fmt.Sprintf("unknown(%d)", t)
	}
}
