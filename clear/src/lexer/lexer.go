package lexer

// Look into using channels for error handling and logging

type Lexer struct {
	// Input
	src string

	// Runtime state
	pos     int
	readPos int
	ch      byte
	peekCh  byte
	line    int
	col     int

	// Output
	tokens []Token
	errors []error
}
