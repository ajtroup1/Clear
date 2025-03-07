package parser

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ajtroup1/clear/ast"
	"github.com/ajtroup1/clear/errors"
	"github.com/ajtroup1/clear/lexer"
	"github.com/ajtroup1/clear/logger"
	"github.com/ajtroup1/clear/token"
	// "github.com/sanity-io/litter"
)

// PRECEDENCE HANDLING
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	POSTFIX     // x++ or x--
	CALL        // myFunction(x)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.PLUS_EQ:  SUM,
	token.MINUS_EQ: SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MULT_EQ:  PRODUCT,
	token.DIV_EQ:   PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.INC:      POSTFIX,
	token.DEC:      POSTFIX,
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	Errors []*errors.Error

	log            *logger.Logger
	debug          bool
	encounterCount int

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

func New(l *lexer.Lexer, log *logger.Logger, debug bool) *Parser {
	p := &Parser{
		l:              l,
		Errors:         []*errors.Error{},
		log:            log,
		debug:          debug,
		encounterCount: 1,
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.PLUS_EQ, p.parseInfixExpression)
	p.registerInfix(token.MINUS_EQ, p.parseInfixExpression)
	p.registerInfix(token.MULT_EQ, p.parseInfixExpression)
	p.registerInfix(token.DIV_EQ, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.INC, p.parsePostfixExpression)
	p.registerPostfix(token.DEC, p.parsePostfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s ('%s') instead",
		t, p.peekToken.Type, p.peekToken.Literal)
	err := errors.Error{
		Message: msg,
		Line:    p.peekToken.Line,
		Col:     p.peekToken.Col,
		Stage:   "Parsing",
		Context: p.l.Lines[p.peekToken.Line-1],
	}
	p.Errors = append(p.Errors, &err)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	if t == token.ILLEGAL {
		// Illegal tokens are reported by the lexer
		return
	}
	if isStatement(t) {
		msg := fmt.Sprintf("'%s' statement not allowed as expression", t)
		err := errors.Error{
			Message: msg,
			Line:    p.curToken.Line,
			Col:     p.curToken.Col,
			Stage:   "Parsing",
			Context: p.l.Lines[p.curToken.Line-1],
		}
		p.Errors = append(p.Errors, &err)
		return
	}

	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	err := errors.Error{
		Message: msg,
		Line:    p.curToken.Line,
		Col:     p.curToken.Col,
		Stage:   "Parsing",
		Context: p.l.Lines[p.curToken.Line-1],
	}
	p.Errors = append(p.Errors, &err)
}

func isStatement(t token.TokenType) bool {
	switch t {
	case token.LET, token.RETURN, token.MOD:
		return true
	}
	return false
}

func (p *Parser) ParseProgram() *ast.Program {
	if p.debug {
		p.log.AppendParser("### Live Encounters:\n\n")
	}
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Starting to parse the program node...\n\n", p.encounterCount))
		p.encounterCount++
		p.log.AppendParser("\t- This requires invoking a loop until end of file is reached, and parsing statements one-by-one until that point. As statements are parsed, they are appended to the `Program`'s `Statements` slice\n\n")
	}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		result := checkNilStmt(stmt)
		if result {
			continue
		} else {
			if p.debug {
				p.log.AppendParser(fmt.Sprintf("%d. Parsed a statement to append to program's `Statements` slice: `%s`\n", p.encounterCount, stmt.String()))
			}
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	// Filtering the statements
	var filteredStatements []ast.Statement
	for _, stmt := range program.Statements {
		if module, ok := stmt.(*ast.ModuleStatement); ok {
			program.Modules = append(program.Modules, module)
		} else {
			filteredStatements = append(filteredStatements, stmt)
		}
	}

	program.Statements = filteredStatements

	if len(program.Statements) == 1 {
		if stmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok && stmt.Expression == nil {
			program.NoStatements = true
		}
	}

	if len(program.Statements) == 0 {
		program.NoStatements = true
	}

	if p.debug {
		toks := p.l.Tokens

		sort.Slice(toks, func(i, j int) bool {
			if toks[i].Line == toks[j].Line {
				return toks[i].Col < toks[j].Col
			}
			return toks[i].Line < toks[j].Line
		})

		parsingDescription := "Parsing is the organization of our tokens into a tree structure that represents the program. This is done by reading the tokens from the lexer sequentially and forming 'nodes' on the fly.\n\n- Nodes are the building blocks of any AST and can range from:\n\t- LetStatement{ Token: Let, Name: 'x', Value: 5 }\n\t- Integer{ Token: INT, Value: 0 }\n\t- InfixExpression{ Token: +, Left: 'x', Operator: '+', Right: '5' }\n\t\n\nSo, essentially these nodes are the next level of abstraction from tokens and represent the structure of the program in a more complex and well-structured manner that is easy to evaluate later.\n\n"

		p.log.Append("\n\n**Here is the stream of all tokens generated by the lexer:**\n\n```\n")
		for _, tok := range toks {
			p.log.Append(fmt.Sprintf("%s\n", tok.String()))
		}
		p.log.Append("```")
		p.log.DefineSection("Parsing", parsingDescription)
		p.log.Append(p.log.GetParserLog())
		jsonData, err := json.MarshalIndent(program, "", "  ")
		if err != nil {
			p.log.Append(fmt.Sprintf("Error serializing program to JSON: %s", err))
			return nil
		}
		p.log.Append(fmt.Sprintf("\n\n**Successfully parsed the program!**\n\nHere is your program node in tree format:\n```\n%s\n```\n", string(jsonData)))
	}

	return program
}

func checkNilStmt(stmt ast.Statement) bool {
	if stmt == nil {
		return true
	}
	if letStmt, ok := stmt.(*ast.LetStatement); ok {
		if letStmt == nil {
			return true
		}
	}
	if retStmt, ok := stmt.(*ast.ReturnStatement); ok {
		if retStmt == nil {
			return true
		}
	}
	if exp, ok := stmt.(*ast.ExpressionStatement); ok {
		if exp == nil {
			return true
		}
	}
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}
