package lexer

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/clear/errors"
	"github.com/ajtroup1/clear/logger"
	"github.com/ajtroup1/clear/token"
)

type Lexer struct {
	input        string
	Lines        []string
	linePos      int
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	col          int

	log            *logger.Logger
	debug          bool
	Tokens         []token.Token
	encounterCount int

	Errors []*errors.Error
}

func New(input string, lo *logger.Logger, debug bool) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0, Errors: []*errors.Error{}, log: lo, debug: debug, encounterCount: 1}
	if debug {
		l.log.DefineSection("Lexical Analysis / Lexing / Tokenization", "Lexing (or tokenization) is the process of converting a sequence of characters into a sequence of tokens.\n\nThese tokens are the simplest level of strutured data pertaining to the source code information.\n\n*Example Token*: `let` (Type: LET, Literal: 'let')\n\n- Optionally, the token can track other information such as line and column information, which is used for error reporting.\n\n\t- Token: `return` (Type: RETURN, Literal:'return', Line: 6, Column: 12)\n\nThe lexer reads the source code character by character and generates tokens based on the characters it reads. The lexer is also the first step in the compilation or interpretation process. The lexer is additionally responsible for removing whitespace and comments from the source code.")
		l.log.Append("**Source code:** \n```js\n")
		l.log.Append(input)
		l.log.Append("\n```\n")
		l.log.Append("\nInitializing lexer...\n")
		l.log.Append("\n**Lexing source code...**\n\n")
		l.log.Append("### Live encounters:\n\n")
	}
	l.readChar()
	l.Lines = strings.Split(input, "\n")
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
		} else {
			tok = l.newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal}
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
		} else {
			tok = l.newToken(token.SLASH, l.ch)
		}
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)
	case '<':
		tok = l.newToken(token.LT, l.ch)
	case '>':
		tok = l.newToken(token.GT, l.ch)
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
			tok.Type = token.LookupIdent(tok.Literal, l.log, l.encounterCount)
			tok.Line = l.line
			tok.Col = l.col - len(tok.Literal)
			l.Tokens = append(l.Tokens, tok)
			l.log.Append(fmt.Sprintf("%d. Tokenized Token::%s '%s' at [line: %d, col: %d]\n", l.encounterCount, tok.Type, tok.Literal, tok.Line, tok.Col))
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
			l.Tokens = append(l.Tokens, tok)
			l.log.Append(fmt.Sprintf("%d. Tokenized Token::%s '%s' at [line: %d, col: %d]\n", l.encounterCount, tok.Type, tok.Literal, tok.Line, tok.Col))
			return tok
		} else {
			err := errors.New("illegal character '"+string(l.ch)+"'", l.line, l.col, "lexer", l.Lines, false)
			l.Errors = append(l.Errors, err)
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	l.Tokens = append(l.Tokens, tok)
	l.log.Append(fmt.Sprintf("%d. Tokenized Token::%s '%s' at [line: %d, col: %d]\n", l.encounterCount, tok.Type, tok.Literal, tok.Line, tok.Col))
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	l.log.Append(fmt.Sprintf("%d. Comment encountered at [line: %d, col: %d]. Skipping over text...\n", l.encounterCount, l.line, l.col))
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
		l.log.Append(fmt.Sprintf("%d. Encountered newline at [line: %d, col: %d]\n", l.encounterCount, l.line, l.col))
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
	for isLetter(l.ch) {
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
