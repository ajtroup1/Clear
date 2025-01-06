/*

	The Lexer converts source code into tokens, which are the building blocks for any AST
	The Parser relies on Tokens to see structured information about the source code and turn it into a hierarchal tree to traverse later
	Lexing is the first and probably easiest step of converting source code into either machine code, assembly, C, or whatever
	This file defines the necessary functionality for returning a stream of Tokens from a given string of source code

*/

package lexer

import (
	"github.com/ajtroup1/goclear/lexing/token"
)

var keywords = map[string]token.TokenType{
	"true":     token.TRUE,
	"false":    token.FALSE,
	"int":      token.INT,
	"float":    token.FLOAT,
	"string":   token.STRING,
	"char":     token.CHAR,
	"bool":     token.BOOL,
	"void":     token.VOID,
	"if":       token.IF,
	"else":     token.ELSE,
	"for":      token.FOR,
	"while":    token.WHILE,
	"fn":       token.FUNCTION,
	"let":      token.LET,
	"const":    token.CONST,
	"return":   token.RETURN,
	"break":    token.BREAK,
	"continue": token.CONTINUE,
	"null":     token.NULL,
	"new":      token.NEW,
	"class":    token.CLASS,
	"this":     token.THIS,
	"super":    token.SUPER,
	"static":   token.STATIC,
	"module":   token.MODULE,
	"mod":      token.MODULE,
}

type Lexer struct {
	Tokens []token.Token

	src     string
	pos     int
	readPos int
	c       byte

	line int
	col  int
}

func New(src string) *Lexer {
	if len(src) == 0 {
		return nil
	}
	l := &Lexer{
		src:  src,
		line: 1,
		col:  0,
	}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch l.c {
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.line, Col: l.col}
	case '=':
		if l.peek() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "==", Line: l.line, Col: l.col}
		} else if l.peek() == '>' {
			l.readChar()
			tok = token.Token{Type: token.ARROW, Literal: "=>", Line: l.line, Col: l.col}
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
		} else if l.peek() == '>' {
			l.readChar()
			tok = token.Token{Type: token.ARROW, Literal: "->", Line: l.line, Col: l.col}
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
	case '.':
		tok = token.Token{Type: token.DOT, Literal: string(l.c), Line: l.line, Col: l.col}
	case '%':
		tok = token.Token{Type: token.MODULUS, Literal: string(l.c), Line: l.line, Col: l.col}	
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
	case '"':
		tok = l.readString()
	case '\'':
		tok = l.readCharacter()
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

func (l *Lexer) readString() token.Token {
	startLine, startCol := l.line, l.col

	l.readChar()
	pos := l.pos
	for l.c != '"' && l.c != 0 {
		l.readChar()
	}
	literal := l.src[pos:l.pos]
	return token.Token{Type: token.STRING, Literal: literal, Line: startLine, Col: startCol}
}

func (l *Lexer) readCharacter() token.Token {
	l.readChar() // The opening single quote

	tok := token.Token{Type: token.CHAR, Literal: string(l.c), Line: l.line, Col: l.col}

	if l.peek() != '\'' {
		panic("Unterminated character literal")
	}
	l.readChar() // The closing single quote

	return tok
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

func (l *Lexer) readChar() {
    if l.readPos >= len(l.src) {
        l.c = 0
    } else {
        l.c = l.src[l.readPos]
    }
    if l.c == '\n' {
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

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func (l *Lexer) skipWhitespace() {
    for {
        for isWhitespace(l.c) {
            l.readChar()
        }

        if l.c == '/' && l.peek() == '/' {
            for l.c != '\n' && l.c != 0 {
                l.readChar()
            }
            continue
        }

        if l.c == '/' && l.peek() == '*' {
            l.readChar() 
            l.readChar() 
            for {
                if l.c == '*' && l.peek() == '/' {
                    l.readChar()
                    l.readChar() 
                    break
                }
                if l.c == 0 {
                    panic("Unterminated multi-line comment")
                }
                l.readChar()
            }
            continue
        }

        break
    }
}

