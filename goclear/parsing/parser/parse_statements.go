// Outsourcing the parsing of statements to a separate file

package parser

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/goclear/lexing/token"
	"github.com/ajtroup1/goclear/parsing/ast"
)

// Core parsing functions for recognizing statement keywords and parsing the corresponding statements
func (p *Parser) parseStatement() ast.Statement {
	// Execute the corresponding parsing function based on the *current* token type
	// fmt.Printf("p.curToken.Type: %s\n", p.curToken.Type)
	switch p.curToken.Type {
	// If a data type is found, it is an assignment statement
	case token.INT, token.FLOAT, token.STRING, token.CHAR, token.BOOL:
		return p.parseAssignStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.MODULE:
		return p.parseModuleStatement()
	case token.EOF:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

// Simply parses a block of statements recursively and returns them in a slice under a BlockStatement node
//
// Example:
//
//	{
//		int x = 5;
//		int y = 10;
//		return x + y;
//	}
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{BaseNode: ast.BaseNode{Token: p.curToken}}
	block.Statements = []ast.Statement{}

	p.nextToken() // Move past the '{'

	// Until the end of the BlockStatement ('}') or the end of the file, parse statements and append them to the block
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// Parses an assignment statement, which can be a declaration with or without an value assignment
// - Value can be nil if no assignment is made ("int x;")
//
// Example:
//
//	int x = 5;
func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	// All assignment statements must start with a data type
	// It is confirmed that the current token is a data type by parseStatement()
	stmt.Type = mapTokenTypeToDataType(p.curToken.Type)

	// The second token must be an identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Assign the identifier to the statement and check if it has already been declared in the current scope
	stmt.Name = &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
	if _, found := p.symbols[stmt.Name.Value]; found {
		p.addError(fmt.Sprintf("variable '%s' already declared", stmt.Name.Value), p.curToken.Line, p.curToken.Col)
	}

	// If the next token is an assignment operator, parse the expression and assign it to the statement
	if p.peekToken.Type == token.ASSIGN {
		p.nextToken()
		p.nextToken()

		stmt.Value = p.parseExpression(LOWEST)

		if stmt.Type != stmt.Value.GetType() {
			p.addError(fmt.Sprintf(
				"type mismatch: cannot assign %s to %s", stmt.Value.GetType(), stmt.Type),
				p.curToken.Line, p.curToken.Col)
		}
	}

	// Optional semicolon at the end of the statement
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// Register the variable in the current scope
	p.symbols[stmt.Name.Value] = Symbol{
		Name:  stmt.Name.Value,
		Type:  stmt.Type,
		Value: stmt.Value,
	}

	return stmt
}

// Parses a reassignment statement, which is a statement that reassigns a value to an already declared variable
func (p *Parser) parseReassignStatement() *ast.ReassignStatement {
	stmt := &ast.ReassignStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	stmt.Name = &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal}
	if _, found := p.symbols[stmt.Name.Value]; !found {
		p.addError(fmt.Sprintf("variable '%s' not declared", stmt.Name.Value), p.curToken.Line, p.curToken.Col)
	}

	p.nextToken() // Consume '='
	p.nextToken() // Move to the expression

	stmt.Value = p.parseExpression(LOWEST)
	stringType := strings.Split(stmt.Value.ToString(), " ")[0]
	if p.symbols[stmt.Name.Value].Type != mapTokenTypeToDataType(token.TokenType(stringType)) {
		p.addError(fmt.Sprintf(
			"type mismatch: cannot assign %s to %s for var '%s'", stmt.Value.GetType(), p.symbols[stmt.Name.Value].Type, stmt.Name.Value),
			p.curToken.Line, p.curToken.Col)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	p.symbols[stmt.Name.Value] = Symbol{
		Name:  stmt.Name.Value,
		Type:  p.symbols[stmt.Name.Value].Type,
		Value: stmt.Value,
	}

	return stmt
}

// Parses a return statement, which is a statement that returns a value from a function
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	stmt.Init = p.parseStatement()

	if !p.curTokenIs(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)
	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	stmt.Post = p.parseExpression(LOWEST)

	p.nextToken()
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseModuleStatement() *ast.ModuleStatement {
    stmt := &ast.ModuleStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

    if !p.expectPeek(token.IDENT) {
        return nil
    }

		ident := &ast.Identifier{BaseNode: ast.BaseNode{Token: p.curToken}, Value: p.curToken.Literal, Type: ast.MODULE}

    stmt.Name = ident

    if !p.expectPeek(token.COLON) {
        return nil
    }

    if p.peekTokenIs(token.ASTERISK) {
        stmt.ImportAll = true
        p.nextToken()
    } else if p.peekTokenIs(token.LBRACKET) {
        stmt.Imports = p.parseModuleImports()
    } else {
        p.addError("expected '*' or list of imports after module statement", p.peekToken.Line, p.peekToken.Col)
        return nil
    }

    if p.peekTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

		fmt.Printf("stmt: %v\n", stmt)

    return stmt
}

func (p *Parser) parseModuleImports() []*ast.Identifier {
    var imports []*ast.Identifier
    p.nextToken() // Skip '['

    for !p.peekTokenIs(token.RBRACKET) {
        if !p.peekTokenIs(token.IDENT) {
            p.addError("expected identifier in imports list", p.peekToken.Line, p.peekToken.Col)
            return nil
        }

				ident := &ast.Identifier{BaseNode: ast.BaseNode{Token: p.peekToken}, Value: p.peekToken.Literal, Type: ast.MODULEFUNCTION}

        imports = append(imports, ident)
        p.nextToken()

        if p.peekTokenIs(token.COMMA) {
            p.nextToken()
        }
    }

		fmt.Printf("imports: %v\n", imports)

    p.nextToken() // Skip ']'
    return imports
}

// Parses an expression and wraps it in an ExpressionStatement node
// This is necessary for expressions that are not part of a larger statement to be appended to the program's statements slice
func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{BaseNode: ast.BaseNode{Token: p.curToken}}

	if p.peekToken.Type == token.ASSIGN {
		return p.parseReassignStatement()
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
