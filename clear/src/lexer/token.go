package lexer

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and literals
	IDENT  // main, foo, bar, etc.
	// Typing tokens
	INT    // 12345
	STRING // "foobar"
	BOOL   // true, false

	// Arithmetic operators
	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	MOD      // %

	// Comparison operators
	EQ     // ==
	NOT_EQ // !=
	LT     // <
	GT     // >
	LTE    // <=
	GTE    // >=

	// Logical operators
	AND // &&
	OR  // ||

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]

	// Keywords
	FN          // fn
	LET         // ? Why would i need let or var? Go is statically typed but it uses var
	CONST       // const
	VAR         // var
	INT_TYPE    // int
	STRING_TYPE // string
	BOOL_TYPE   // bool
	TRUE        // true
	FALSE       // false
	IF          // if
	ELSE        // else
	RETURN      // return
)
