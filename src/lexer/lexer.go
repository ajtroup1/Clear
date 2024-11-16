/*
	The Lexer reads input src code char by char and forms a stream tokens accordingly
	Tokens are little building blocks for bigger blocks of instructions created in the parser
		Tokens contain segments such as: "let", "x", "=", "5", ";"
	Probably the easiest and least flexible component of any programming language
*/

package lexer

import "github.com/ajtroup1/clear/src/token"

// Struct for the Lexer and its state
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// Returns a new instance of Lexer given an input src
// Also starts lexing by initially calling readChar()
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Essential lexing function that reads the tracked char and advances the lexer state accordingly
// Also handles for end of input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 'Peeks" ahead to return the tracked char but does not advance the lexer state
// Used for conditions, deciding whether or not to advance, etc...
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

/*
	Core of the Lexer process
	Reads an individual char and switches it, assigning it accordingly and returning the token
	Some assignments require multiple chars to be read
		These are buffered accordingly, such as the default for the switch
*/

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// Skip any whitespace to isolate the char
	l.skipWhitespace()

	// Switch and evaluate what token to return
	switch l.ch {
	case '=':
		// Could either be "=" or "=="
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		// Could either be "!" or "!="
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// Wasnt a specific char, so must be an identifier or number
		// Can either be a user-defined identifier or a Clear reserved keyword
		if isLetter(l.ch) {
			// Idents and keywords must start with a letter
			// Or else, they have to be evaluated to a number
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// Helper function to abstract creating and returning a new token
// **Given a token type and literal value**
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Buffers the input until the lexer no longer encounters a letter or digit and returns it
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Returns a bool indicating whether the tracked char is a alphabetical character
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// You can guess what this one does
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Clear does not consider whitespace and skips any whitespace characters before reading the next token
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// Buffers the input until a non-number is reached and returns it
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
