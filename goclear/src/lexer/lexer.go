package lexer

import "github.com/ajtroup1/goclear/src/token"

type Lexer struct {
	input        string
	position     int      // current position in input (points to current char)
	readPosition int      // current reading position in input (after current char)
	ch           byte     // current char under examination
	line         int      // current line number
	column       int      // current column number
	lines        []string // stores all input lines
	currentLine  []byte   // stores current line content being read
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.lines = append(l.lines, "") 
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(l, token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(l, token.PLUS, l.ch)
	case '-':
		tok = newToken(l, token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(l, token.BANG, l.ch)
		}
	case '/':
		tok = newToken(l, token.SLASH, l.ch)
	case '*':
		tok = newToken(l, token.ASTERISK, l.ch)
	case '<':
		tok = newToken(l, token.LT, l.ch)
	case '>':
		tok = newToken(l, token.GT, l.ch)
	case ';':
		tok = newToken(l, token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(l, token.COMMA, l.ch)
	case '{':
		tok = newToken(l, token.LBRACE, l.ch)
	case '}':
		tok = newToken(l, token.RBRACE, l.ch)
	case '(':
		tok = newToken(l, token.LPAREN, l.ch)
	case ')':
		tok = newToken(l, token.RPAREN, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(l, token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	if l.ch == '\n' {
		l.lines = append(l.lines, string(l.currentLine))
		l.currentLine = []byte{} 
		l.line++
		l.column = 0
	} else {
		l.currentLine = append(l.currentLine, l.ch) 
	}

	l.column++
	l.position = l.readPosition
	l.readPosition += 1
}


func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) LineContent(lineNumber int) string {
	if lineNumber >= 1 && lineNumber <= len(l.lines) {
		return l.lines[lineNumber-1]
	}
	return ""
}

func newToken(l *Lexer, tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    l.line,
		Column:  l.column,
	}
}
