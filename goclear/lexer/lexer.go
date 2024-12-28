package lexer

import "github.com/ajtroup1/goclear/token"

var keywords = map[string]token.TokenType{
	"true":     token.BOOL,
	"false":    token.BOOL,
	"if":       token.IF,
	"else":     token.ELSE,
	"for":      token.FOR,
	"while":    token.WHILE,
	"fn":       token.FUNCTION,
	"let":      token.LET,
	"return":   token.RETURN,
	"break":    token.BREAK,
	"continue": token.CONTINUE,
	"null":     token.NULL,
	"new":      token.NEW,
	"class":    token.CLASS,
	"this":     token.THIS,
	"super":    token.SUPER,
	"static":   token.STATIC,
	"import":   token.IMPORT,
	"export":   token.EXPORT,
}

type Lexer struct {
	src     string
	pos     int
	readPos int
	c       byte

	line int
	col  int
}

func New(src string) *Lexer {
	l := &Lexer{
		src:  src,
		line: 1,
		col:  0,
	}
	l.readChar()
	return l
}

func (l *Lexer) Lex() []token.Token {
	var tokens []token.Token
	for l.c != 0 {
		tok := l.nextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) nextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch l.c {
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.line, Col: l.col}
	case '=':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "==", Line: l.line, Col: l.col}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	case '+':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.PLUS_EQ, Literal: "+=", Line: l.line, Col: l.col}
		} else if l.peek() == '+' {
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: "++", Line: l.line, Col: l.col}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	case '-':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.MINUS_EQ, Literal: "-=", Line: l.line, Col: l.col}
		} else if l.peek() == '-' {
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: "--", Line: l.line, Col: l.col}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	case '!':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!=", Line: l.line, Col: l.col}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	case '*':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.MUL_EQ, Literal: "*=", Line: l.line, Col: l.col}
		} else {
			tok = token.Token{Type: token.ASTERISK, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	case '/':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.DIV_EQ, Literal: "/=", Line: l.line, Col: l.col}
		} else {
			tok = token.Token{Type: token.SLASH, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	case '<':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: "<=", Line: l.line, Col: l.col}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	case '>':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: ">=", Line: l.line, Col: l.col}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.c), Line: l.line, Col: l.col}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.c), Line: l.line, Col: l.col}
	case ':':
		tok = token.Token{Type: token.COLON, Literal: string(l.c), Line: l.line, Col: l.col}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.c), Line: l.line, Col: l.col}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.c), Line: l.line, Col: l.col}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.c), Line: l.line, Col: l.col}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.c), Line: l.line, Col: l.col}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.c), Line: l.line, Col: l.col}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.c), Line: l.line, Col: l.col}
	default:
		if isLetter(l.c) {
			return l.readIdentifier()
		} else if isDigit(l.c) {
			return l.readNumber()
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.c), Line: l.line, Col: l.col}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() token.Token {
	startLine, startCol := l.line, l.col

	pos := l.pos
	for isLetter(l.c) {
		l.readChar()
	}
	literal := l.src[pos:l.pos]
	return keywordLookup(literal, startLine, startCol)
}

func keywordLookup(ident string, line, col int) token.Token {
	if tok, found := keywords[ident]; found {
		return token.Token{Type: tok, Literal: ident, Line: line, Col: col}
	}
	return token.Token{Type: token.IDENT, Literal: ident, Line: line, Col: col}
}

func (l *Lexer) readNumber() token.Token {
	startLine, startCol := l.line, l.col

	pos := l.pos
	isFloat := false

	for isDigit(l.c) || l.c == '.' {
		if l.c == '.' {
			isFloat = true
		}
		l.readChar()
	}

	literal := l.src[pos:l.pos]

	if isFloat {
		return token.Token{Type: token.FLOAT, Literal: literal, Line: startLine, Col: startCol}
	}
	return token.Token{Type: token.INT, Literal: literal, Line: startLine, Col: startCol}
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.c) {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.src) {
		l.c = 0
	} else {
		l.c = l.src[l.readPos]
	}

	if l.c == '\n' || l.c == '\r' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}

	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.src) {
		return 0
	}

	return l.src[l.readPos]
}

func (l *Lexer) peekN(n int) byte {
	if l.readPos+n >= len(l.src) {
		return 0
	}

	return l.src[l.readPos+n]
}

func (l *Lexer) consume() {
	if l.readPos >= len(l.src) {
		l.c = 0
	} else {
		l.c = l.src[l.readPos]
	}
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
