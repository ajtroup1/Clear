package lexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ajtroup1/goclear/src/util"
)

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

// Overall lexer class to manage information and state
type Lexer struct {
	Patterns []regexPattern // List of regex patterns and their corresponding handlers
	Tokens   []Token        // List of tokens generated from lexing
	src      string         // Input source code passed as a string
	current  int            // Current lexer position in the source code

	// Error handling
	line int // Current line number
	col  int // Current column number
}

// Lexer constructor
// Takes in an entire source as a string casted from a byte array
func NewLexer(src string) *Lexer {
	return &Lexer{
		src:     src,
		current: 0,
		line:    1,
		col:     1,
		Tokens:  make([]Token, 0),
		Patterns: []regexPattern{
			{regexp.MustCompile(`\s+`), whitespaceHandler},
			{regexp.MustCompile(`\/\/.*`), whitespaceHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_BRACE, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_BRACE, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(COMPARISON, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUAL, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(BANG, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUAL, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUAL, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUAL, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUAL, "-=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(MINUS, "-")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

// This handler can handle most simple patterns in Clear
// Returns a simple function that advances the length of the token and pushes a new token to the Lexer.Token slice
func defaultHandler(t TokenType, lit string) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		lex.consumeN(len(lit))
		lex.push(NewToken(t, lit))
	}
}

// Special regex handlers
func numberHandler(l *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(l.peekRemainder())
	l.push(NewToken(NUMBER, match))
	l.consumeN(len(match))
}

func stringHandler(l *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(l.peekRemainder())
	stringLit := l.peekRemainder()[match[0]+1 : match[1]-1]

	l.push(NewToken(STRING, stringLit))
	l.consumeN(len(stringLit) + 2)
}

func symbolHandler(l *Lexer, regex *regexp.Regexp) {
	val := regex.FindString(l.peekRemainder())

	if t, exists := keyword_lookup[val]; exists {
		l.push(NewToken(t, val))
	} else {
		l.push(NewToken(IDENT, val))
	}

	l.consumeN(len(val))
}

func whitespaceHandler(l *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(l.peekRemainder())
	l.consumeN(match[1])
}

// Lexer state helpers

func (l *Lexer) peek() byte {
	return l.src[l.current]
}

func (l *Lexer) peekRemainder() string {
	return l.src[l.current:]
}

func (l *Lexer) consumeN(n int) {
	for i := 0; i < n; i++ {
		if l.src[l.current] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		l.current++
	}
}

func (l *Lexer) atEOF() bool {
	return l.current >= len(l.src)
}

func (l *Lexer) push(t Token) {
	l.Tokens = append(l.Tokens, t)
}

// Core lexer function
// Advances through the entire source and creates according tokens in a stream
func Tokenize(src string) []Token {
	// Instantiate a new lexer with the given source code
	lex := NewLexer(src)

	// Until the lexer reaches the end of the file, create tokens and append them to the stream
	for !lex.atEOF() {
		matched := false // Used later to check if a token is unrecognized and print an error
		for _, pattern := range lex.Patterns {
			loc := pattern.regex.FindStringIndex(lex.peekRemainder())

			// The location must be found (!nil) AND have matched at the current position (ex. not 5 positions ahead)
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			lineContent := getLineContent(lex.src, lex.line)
			util.PrintErrorPanic(
				"Lexer",
				fmt.Sprintf(
					"unrecognized token '%s' on line %d, col %d\n%s\n%s^",
					lex.peekRemainder()[0:1], lex.line, lex.col, lineContent, strings.Repeat(" ", lex.col-1),
				),
			)
		}
	}

	lex.push(NewToken(EOF, "EOF"))
	return lex.Tokens
}

// Returns the entire line given the source and line number to print
func getLineContent(src string, line int) string {
	lines := strings.Split(src, "\n")
	if line-1 < len(lines) {
		return lines[line-1]
	}
	return ""
}
