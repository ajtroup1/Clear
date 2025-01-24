package lexer

import (
	"fmt"

	"github.com/ajtroup1/compiled_clear/src/errorlogger"
	"github.com/ajtroup1/compiled_clear/src/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	col          int

	el    *errorlogger.ErrorLogger
	debug bool
}

func New(input string, el *errorlogger.ErrorLogger, debug bool) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0, el: el, debug: debug}

	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUS_EQ, Literal: literal}
		} else {
			tok = l.newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MINUS_EQ, Literal: literal}
		} else {
			tok = l.newToken(token.MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = l.newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DIV_EQ, Literal: literal}
		} else {
			tok = l.newToken(token.SLASH, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MULT_EQ, Literal: literal}
		} else {
			tok = l.newToken(token.ASTERISK, l.ch)
		}
	case '<':
		tok = l.newToken(token.LT, l.ch)
	case '>':
		tok = l.newToken(token.GT, l.ch)
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal}
		} else {
			l.el.NewError(l.line, l.col, fmt.Sprintf("Illegal character '%s' encountered", string(l.ch)), "Lexer", false)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal}
		} else {
			l.el.NewError(l.line, l.col, fmt.Sprintf("Illegal character '%s' encountered", string(l.ch)), "Lexer", false)
		}
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = l.newToken(token.COLON, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '.':
		tok = l.newToken(token.DOT, l.ch)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Col = l.col - len(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			lit, isFLoat := l.readNumber()
			if isFLoat {
				tok.Type = token.FLOAT
			}
			tok.Literal = lit
			tok.Line = l.line
			tok.Col = l.col - len(tok.Literal)
			return tok
		} else {
			l.el.NewError(l.line, l.col, fmt.Sprintf("Illegal character '%s' encountered", string(l.ch)), "Lexer", false)
			tok = l.newToken(token.ILLEGAL, l.ch)
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

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
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
		l.line++
		l.col = 0
	} else {
		l.col++
	}

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
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, bool) {
	isFloat := false
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' {
		isFloat = true
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position], isFloat
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: l.line, Col: l.col}
}
