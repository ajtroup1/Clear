package lexer

import (
	"strings"
	"testing"

	"github.com/ajtroup1/compiled_clear/src/errorlogger"
	"github.com/ajtroup1/compiled_clear/src/token"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []token.Token
	}{
		{
			`=+(){},;`,
			[]token.Token{
				{Type: token.ASSIGN, Literal: "=", Line: 1, Col: 1},
				{Type: token.PLUS, Literal: "+", Line: 1, Col: 2},
				{Type: token.LPAREN, Literal: "(", Line: 1, Col: 3},
				{Type: token.RPAREN, Literal: ")", Line: 1, Col: 4},
				{Type: token.LBRACE, Literal: "{", Line: 1, Col: 5},
			},
		},
		{
			`let five = 5;`,
			[]token.Token{
				{Type: token.LET, Literal: "let", Line: 1, Col: 1},
				{Type: token.IDENT, Literal: "five", Line: 1, Col: 5},
				{Type: token.ASSIGN, Literal: "=", Line: 1, Col: 10},
				{Type: token.INT, Literal: "5", Line: 1, Col: 12},
				{Type: token.SEMICOLON, Literal: ";", Line: 1, Col: 13},
			},
		},
		{
			`let add = fn(x, y) {
				x + y;
				};`,
			[]token.Token{
				{Type: token.LET, Literal: "let", Line: 1, Col: 1},
				{Type: token.IDENT, Literal: "add", Line: 1, Col: 5},
				{Type: token.ASSIGN, Literal: "=", Line: 1, Col: 9},
				{Type: token.FUNCTION, Literal: "fn", Line: 1, Col: 11},
				{Type: token.LPAREN, Literal: "(", Line: 1, Col: 13},
				{Type: token.IDENT, Literal: "x", Line: 1, Col: 14},
				{Type: token.COMMA, Literal: ",", Line: 1, Col: 15},
				{Type: token.IDENT, Literal: "y", Line: 1, Col: 17},
				{Type: token.RPAREN, Literal: ")", Line: 1, Col: 18},
				{Type: token.LBRACE, Literal: "{", Line: 1, Col: 20},
				{Type: token.IDENT, Literal: "x", Line: 2, Col: 5},
				{Type: token.PLUS, Literal: "+", Line: 2, Col: 7},
				{Type: token.IDENT, Literal: "y", Line: 2, Col: 9},
				{Type: token.SEMICOLON, Literal: ";", Line: 2, Col: 10},
				{Type: token.RBRACE, Literal: "}", Line: 3, Col: 5},
				{Type: token.SEMICOLON, Literal: ";", Line: 3, Col: 6},
			},
		},
	}

	for _, tt := range tests {
		el := errorlogger.New(strings.Split(tt.input, "\n"), false)
		l := New(tt.input, el, false)

		for i, expected := range tt.expected {
			tok := l.NextToken()

			if tok.Type != expected.Type {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, expected.Type, tok.Type)
			}

			if tok.Literal != expected.Literal {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
					i, expected.Literal, tok.Literal)
			}

			if tok.Line != expected.Line {
				t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
					i, expected.Line, tok.Line)
			}

			if tok.Col != expected.Col {
				t.Fatalf("tests[%d] - col wrong. expected=%d, got=%d (%s)",
					i, expected.Col, tok.Col, tok.Literal)
			}
		}
	}
}
