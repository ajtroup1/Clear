package parser

import (
	"fmt"
	"strconv"

	"github.com/ajtroup1/goclear/src/ast"
	"github.com/ajtroup1/goclear/src/lexer"
	"github.com/ajtroup1/goclear/src/token"
	"github.com/ajtroup1/goclear/src/types"
)

// Precedence enum for operator precedence and binary expression parsing
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// Types of parsing functoins
// - Prefix: "!myVar", "-10"
// - Infix: "1 + 2", "myObject.property"
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// Parser struct tracks state, errors, and handlers for parsing
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []types.ParserError

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// Register a prefix parsing function given a token type
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Register an infix parsing function given a token type
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: make([]types.ParserError, 0),
	}

	p.nextToken()
	p.nextToken()

	// Register prefix functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	// Since grouped expressions are handled with just the left paren above, maybe add a handler for right paren that returns a parser error?
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// Register infix functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

// Simply return the Parser's error
func (p *Parser) Errors() []types.ParserError {
	return p.errors
}

// Add an error if the peeked token is not an expected type
func (p *Parser) peekError(t token.Token) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t.Type, p.peekToken.Type)
	pe := types.ParserError{
		Message:     msg,
		Line:        t.Line,
		Column:      t.Column,
		LineContent: p.l.LineContent(t.Line),
	}
	p.errors = append(p.errors, pe)
}

// Advance the Lexer state 'by 1'
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Highest level function to parse the entire program into a Statement[] in an ast.Program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// Calls parsing for any statement according to the current token type
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// -----------------
// STATEMENT PARSING
// -----------------

func (p *Parser) parseLetStatement() *ast.LetStatement {
    // * Example parse: "[let] [x] [=] [7] [;]"
    // * "let x = 7;"
    // Need to assign:
    // 	1 - The LET token
    // 	2 - The identifier / Name property (basically just a string)
    // 	3 - The value assigned to the variable (can be any expression)

    // [let] x = 7;
    //   ^
    stmt := &ast.LetStatement{Token: p.curToken}
    if !p.expectPeek(token.IDENT) {
        return nil
    }

    // let [x] = 7;
    //      ^
    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // Identifier struct simply wraps the literal name value with the IDENT token
    // let x [=] 7;
    //        ^
    if !p.expectPeek(token.ASSIGN) {
        return nil
    }

    // let x = [7];
    //          ^
    p.nextToken()
    stmt.Value = p.parseExpression(LOWEST)

    // let x = 7[;]
    //           ^
    for !p.curTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    // Now, we should have:
    // ast.LetStatement {
    // 	Token: LET,
    // 	Name: ast.Identifier : Expression { Token: IDENT, Value: "x" },
    // 	Value: ast.IntegerLiteral : Expression {
    // 		Token: INT,
    // 		Value: 7,
    // 	},
    // }
    return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// * Example parse: "[return] [0] [;]"
	// * "return 0;"
	// We just need to assign the ReturnValue, which can be done
	// by parsing the expression next to the keyword

	// [return] 0;
	//    ^
	stmt := &ast.ReturnStatement{Token: p.curToken} // The RETURN token
	p.nextToken()

	// return [0];
	//         ^
	// The '0' is really an expression
	// '0' could really be replaced with any expression, such as:
	// 	- "x + 7"
	// 	- "myObject.property"
	// 	- "myFunction()"
	stmt.ReturnValue = p.parseExpression(LOWEST)

	// return 0 [;]
	//          ^
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// Now, we should have:
	// ast.ReturnStatement {
	// 	Token: RETURN,
	// 	ReturnValue: ast.IntegerLiteral : Expression {
	// 		Token: INT,
	// 		Value: 0,
	// 	},
	// }
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// * Example parse: "[y] [*] [2] [;]"
	// * "y * 2;"
	// Note*: this expression should be the only code on that line
	// Expression statements are just a way to append a lone statement to the Program.Statement[]

	// [y] * 2;
	//  ^
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	// [y] [*] [2] ;
	// [^   ^   ^]
	// parseExpression() will evaluate this as an InfixExpression
	// In this example, it will return:
	// ast.InfixExpression : Expression {
	// 	Left: ast.Identifier : Expression { Token: IDENT, Value: "y" },
	// 	Operator: "*",
	// 	Right: ast.IntegerLiteral : Expression {
	// 		Token: INT,
	// 		Value: 2,
	// 	},
	// }
	stmt.Expression = p.parseExpression(LOWEST)

	// y + 2 [;]
	//        ^
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// Now, we should have:
	// ast.ExpressionStatement : Statement {
	// 	Token: token.Token {
	// 		TokenType: IDENT,
	// 		Literal: "y",
	// 	},
	// 	Expression: ast.InfixExpression : Expression {
	// 		Left: ast.Identifier : Expression { Token: IDENT, Value: "y" },
	// 		Operator: "*",
	// 		Right: ast.IntegerLiteral : Expression {
	// 			Token: INT,
	// 			Value: 2,
	// 		},
	// 	}
	// }
	return stmt
}

// Handle the parser encountering a token with no associated prefixParseFn()
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	pe := types.ParserError{
		Message:     msg,
		Line:        p.curToken.Line,
		Column:      p.curToken.Column,
		LineContent: p.l.LineContent(p.curToken.Line),
	}
	p.errors = append(p.errors, pe)
}

// ------------------
// EXPRESSION PARSING
// ------------------

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// Simple function to parse an identifier:
// x, foo, ...
func (p *Parser) parseIdentifier() ast.Expression {
	// Should return something like:
	// ast.Identifier { Token: IDENT, Value: "x" }
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Handles parsing any integer literal:
// 5, 20, 999999, ...
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// Instantiate the IntegerLiteral instance and assign the INT token to Token
	lit := &ast.IntegerLiteral{Token: p.curToken}
	// Use strconv native package to convert the string literal to and int64
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	// If strconv.ParseInt() failes, append an internal error
	if err != nil {
		msg := fmt.Sprintf("(internal 'go/strconv' error) could not parse %q as integer", p.curToken.Literal)
		pe := types.ParserError{
			Message:     msg,
			Line:        p.curToken.Line,
			Column:      p.curToken.Column,
			LineContent: p.l.LineContent(p.curToken.Line),
		}
		p.errors = append(p.errors, pe)
		return nil
	}

	// Assign the newly parsed int64 to the IntergerLiteral.Value field (int64)
	lit.Value = value

	// Now, we should have something like this:
	// ast.IntegerLiteral : Expression {
	// 	Token: INT,
	// 	Value: 5,
	// }
	return lit
}

// Handles parsing prefix expressions such as:
// -5, !nil, !myObj, ...
func (p *Parser) parsePrefixExpression() ast.Expression {
	// * Example parse: "[-] [x] [;]"
	// * "-x;"

	// Instantiate the PrefixExpression using the current token
	// The current token should be something like '-', '!', ...

	// [-] x;
	//  ^
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// - [x];
	//    ^
	p.nextToken()
	// Evaulate the expression and return the proper Expression instance
	// In this instance we should receive:
	// ast.Identifier : Expression { Token: IDENT, Value: "x" }
	expression.Right = p.parseExpression(PREFIX)

	// Now, we should have:
	// ast.PrefixExpression : Expression {
	// 	Token: MINUS,
	// 	Operator: "-",
	// 	Right: ast.Identifier : Expression { Literal: IDENT, Value: "x" },
	// }
	return expression
}

// Handles Infix / Binary Expressions that contain a left and right hand side
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// * Example parse: "[x] [+] [7] [;]"
	// * "x + 7;"
	// This example is much more naive than what's possible and common with programming
	// You could have to enact a function, access a member of an object, chain multiple infix expressions, ...

	// Instantiate the InfixExpression
	// Left is already passed in, so half of our expression is accounted for and the
	// parser state now rests on the operator in between left and right

	// x [+] 7;
	//    ^
	expression := &ast.InfixExpression{
		Token:    p.curToken,         // PLUS
		Operator: p.curToken.Literal, // "+"
		Left:     left,               // ast.Identifier : Expression { Token: IDENT, Value: "x" }
	}

	// Returns the SUM (4) precedence associated with '+'
	precedence := p.curPrecedence()
	// Advance past the '+'
	p.nextToken()
	// x + [7] ;
	//      ^
	// In this case, parseExpression will parse the '7' wit precedence SUM (4) and return:
	// ast.IntegerLiteral : Expression { Token: INT, Value: 7 }
	expression.Right = p.parseExpression(precedence)

	// Now, we should have:
	// ast.InfixExpression : Expression {
	// 	Token: PLUS,
	// 	Operator: "+",
	// 	Left: ast.Identifier : Expression { Token: IDENT, Value: "x" },
	// 	Right: ast.IntegerLiteral : Expression { Token: INT, Value: 7 },
	// }
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// * Example parse: "[(] [x] [+] [5] [)] [;]"
	// * "(x + 5);"
	// "(x + 5);"

	// [(] x + 5);
	// ^
	p.nextToken()

	// ( [x] [+] [5] )
	//  [ ^   ^   ^ ]
	// Should parse the entire expression "x + 5"
	exp := p.parseExpression(LOWEST)

	// (x + 5 [)]
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Now, we should have:
	// ast.InfixExpression : Expression {
	// 	Left: ast.Identifier { Token: IDENT, Value: x },
	// 	Operator: "+",
	// 	Right: ast.IntegerLiteral : Expression {
	// 		Token: INT,
	// 		Value: 5,
	// 	},
	// }
	return exp
}

// Handles parsing if (and optionally else) expressions
// If's are expressions because they can actually return a true / false value
func (p *Parser) parseIfExpression() ast.Expression {
	// * Example parse "if (x > 5) { return x } else { return y };"
	// * [if] [(] [x] [>] [5] [)] [{] [return] [x] [}] [else] [{] [return] [y] [}]

	// Instantiate that IfExpression instance and assign the IF token to the Token field
	expression := &ast.IfExpression{Token: p.curToken}
	// After the if keyword, the first token should be an opening parenthesis
	// Also will advance past the '(' if expectPeek() succeeds
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	// if ([x] [>] [5]) { return x } else { return y }
	//    [ ^   ^   ^ ]
	// Parses the entire expression within the parentheses
	// Should return an InfixExpression like so:
	// ast.InfixExpression : Expression {
	// 	Token: GT,
	// 	Operator: ">",
	// 	Left: ast.Identifier : Expression { Token: IDENT, Value: "x" },
	// 	Right: ast.IntegerLiteral : Expression { Token: INT, Value: 5 },
	// }
	expression.Condition = p.parseExpression(LOWEST)

	// The next 2 tokens should be a closing parenthesis and open brace respectively
	// Therefor ending the condition and beginning the consequence
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// The consequence is a BlockStatement, meaning an entire script could run as a result
	// In this case, it will return a simple:
	// ast.BlockStatement : Statement {
	// 	Token: token.Token { TokenType: LBRACE, Literal: '{' }
	// 	Statements: [
	// 		ast.ReturnStatement {
	// 			Token: RETURN,
	// 			ReturnValue: ast.IntegerLiteral : Expression {
	// 				Token: INT,
	// 				Value: ast.Identifier : Expression { Token: IDENT, Value: "x" },
	// 			},
	// 		},
	// 	]
	// }
	expression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		// Alternative will do the same thing as consequence, just with "return y" instead of "return x"
		expression.Alternative = p.parseBlockStatement()
	}

	// Now, we should have this:
	// ast.IfExpression : Expression {
	// 	Token: IF,
	// 	Condition: ast.InfixExpression : Expression {
	// 		Token: GT,
	// 		Operator: ">",
	// 		Left: ast.Identifier : Expression { Token: IDENT, Value: "x" },
	// 		Right: ast.IntegerLiteral : Expression { Token: INT, Value: 5 },
	// 	},
	// 	Consequence: ast.BlockStatement : Statement {
	// 		Token: token.Token { TokenType: LBRACE, Literal: '{' }
	// 		Statements: [
	// 			ast.ReturnStatement {
	// 				Token: RETURN,
	// 				ReturnValue: ast.IntegerLiteral : Expression {
	// 					Token: INT,
	// 					Value: ast.Identifier : Expression { Token: IDENT, Value: "x" },
	// 				},
	// 			},
	// 		]
	// 	},
	// 	Alternative: ast.BlockStatement : Statement {
	// 		Token: token.Token { TokenType: LBRACE, Literal: '{' }
	// 		Statements: [
	// 			ast.ReturnStatement {
	// 				Token: RETURN,
	// 				ReturnValue: ast.IntegerLiteral : Expression {
	// 					Token: INT,
	//					Value: ast.Identifier : Expression { Token: IDENT, Value: "x" },
	// 				},
	// 			},
	// 		]
	// 	},
	// }
	return expression
}

// Parses a list of statements into a BlockStatement
// Useful for functions and if/else statements
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// Instantiate new BlockStatement with '{' token
	block := &ast.BlockStatement{Token: p.curToken}
	// New slice for Statements property to append to
	block.Statements = []ast.Statement{}
	p.nextToken()
	// Parse statements until a '}' is encountered, indicating the end of the scope
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		// Parse the statements and append it to the Statements slice
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// Parses an entire function, which has parameters and a block statement
func (p *Parser) parseFunctionLiteral() ast.Expression {
	// We need to assign:
	// 	-	1. All parameters defined for the function
	// 	- 2. The BlockStatement containing functionality
	// Both will use helper functions

	// Instantiate the function with the FN token
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Start reading parameters, indicated by opening parentheses
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Assign parameters to the function, returned from the helper
	lit.Parameters = p.parseFunctionParameters()

	// Indicate the beginning of the function body with a LBRACE
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Assign the function body by simply parsing a block statement
	lit.Body = p.parseBlockStatement()

	// Now, you should have something like:
	// ast.FunctionLiteral : Expression {
	// 	Token: FN,
	// 	Parameters: [
	// 		ast.Identifier : Expression { Token: IDENT, Value: "myParameter" },
	// 		ast.Identifier : Expression { Token: IDENT, Value: "mySecondParam" },
	// 	],
	// 	Body: ast.BlockStatement : Expression {
	// 		Token: LBRACE,
	// 		Statements: [
	// 			ast.ExpressionStatement : Statement {
	// 				Token: IDENT,
	// 				Expression: ast.InfixExpression : Expression {
	// 					Token: PLUS,
	// 					Operator: "+",
	// 					Left: ast.Identifier : Expression { Token: IDENT, Value: "x" },
	// 					Right: ast.Identifier : Expression { Token: IDENT, Value: "y" },
	// 				},
	// 			},
	// 			ast.ReturnStatement : Statement {
	// 				Token: RETURN,
	// 				ReturnValue: ast.Identifier : Expression { Token: IDENT, Value: "x" },
	// 			},
	// 		]
	// 	},
	// }
	return lit
}

// The helper function for returning parameters for functions
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	// Initialize as an empty Identifier list since there can be >= 0 params
	identifiers := []*ast.Identifier{}

	// If there are no parameters, skip the right parenthesis and return the empty slice
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()

	// At this point, there is at least one parameter
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	// As long as we continue to get comma separated params, keep parsing params
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	// No longer encountering comma separated params, we should see a right parenthesis now
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Now we should be returning something like:
	// []*ast.Identifier : Expression [
	// 	ast.Identifier : Expression { Token: IDENT, Value: "firstParameter" },
	// 	ast.Identifier : Expression { Token: IDENT, Value: "secondParameter" },
	// ]
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
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
		p.peekError(p.curToken)
		return false
	}
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
