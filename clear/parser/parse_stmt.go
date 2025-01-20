package parser

import (
	"fmt"

	"github.com/ajtroup1/clear/ast"
	"github.com/ajtroup1/clear/errors"
	"github.com/ajtroup1/clear/token"
)

func (p *Parser) parseStatement() ast.Statement {
	if p.debug && isStatement(p.curToken.Type) {
		p.log.AppendParser(fmt.Sprintf("%d. Encountered a `%s` token, calling `parse%sStatement()`...\n", p.encounterCount, p.curToken.Type, errors.Capitalize(string(p.curToken.Type))))
	}
	switch p.curToken.Type {
	case token.MOD:
		return p.parseModuleStatement()
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	default:
		if p.debug {
			p.log.AppendParser(fmt.Sprintf("%d. Encountered token (`%s`, type '%s') that doesn't have a predefined statement parse function, so it's either an expression or an assignment statement\n", p.encounterCount, p.curToken.Literal, p.curToken.Type))
		}
		return p.parseExpressionOrAssignStatement()
	}
}

func (p *Parser) parseModuleStatement() *ast.ModuleStatement {
	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Steps in parsing module / import statement:\n", p.encounterCount))
		p.log.AppendParser(fmt.Sprintf("\n\ta. Assigning token to the statement to track positioning [line: %d, col: %d]\n", p.curToken.Line, p.curToken.Col))
	}
	stmt := &ast.ModuleStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		msg := fmt.Sprintf("Expected next token to be IDENT, got %s instead", p.peekToken.Type)
		err := errors.New(msg, p.peekToken.Line, p.peekToken.Col, "Parser", p.peekToken.Literal, false)
		p.Errors = append(p.Errors, err)
		return nil
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\tb. Assigning module name `%s` to the import statement\n", p.curToken.Literal))
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.COLON) {
		// TODO: Error handling
		return nil
	}

	if p.peekTokenIs(token.ASTERISK) {
		if p.debug {
			p.log.AppendParser("\n\tc. Encountered a `*` token, signifying this is a wildcard import. Wilcard import simply means **import all** functionality from that module with no need to specify individual functions from the module. Since this is a wildcard import, that also means we have all the information we need, since the module identifier and the fact that we a targeting everything within it is enough to continue.\n")
		}
		stmt.ImportAll = true
		p.nextToken()
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		if p.debug {
			p.log.AppendParser(fmt.Sprintf("%d. Successfully parsed module statement\n", p.encounterCount))
			p.log.AppendParser(fmt.Sprintf("\t- Module name: `%s`\n\t- ", stmt.Name.Value))
			p.log.AppendParser("Imports: ALL ('*')\n")
		}
		return stmt
	}

	if p.debug {
		p.log.AppendParser("\n\tc. **No wildcard import found**, so we're expecting an array of comma-delimited imports\n")
	}

	if !p.expectPeek(token.LBRACKET) {
		// TODO: Error handling
		return nil
	}

	if p.peekTokenIs(token.RBRACKET) {
		// TODO: Warning for empty import list
		if p.debug {
			p.log.AppendParser("\n\tc.1. *Encountered an empty import list, why did you do that?*\n")
		}
		return stmt
	}

	if p.debug {
		p.log.AppendParser("\n\td. Invoking a loop to continue parsing import identifiers until a `]` token is reached, which signifies the end of the import list.\n")
	}
	importCount := 0
	for !p.peekTokenIs(token.RBRACKET) {
		importCount++
		if p.debug {
			p.log.AppendParser(fmt.Sprintf("\n\td.%d. Parsing import identifier `%d`\n", p.encounterCount, importCount))
		}

		if !p.expectPeek(token.IDENT) {
			msg := fmt.Sprintf("Expected next token to be IDENT, got %s instead", p.peekToken.Type)
			err := errors.New(msg, p.peekToken.Line, p.peekToken.Col, "Parser", p.peekToken.Literal, false)
			p.Errors = append(p.Errors, err)
			return nil
		}

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if p.debug {
			p.log.AppendParser(fmt.Sprintf("\n\t- Encountered valid identifier, appending to the list: `%d: %s.%s`\n", importCount, stmt.Name.Value, ident.Value))
		}

		stmt.Imports = append(stmt.Imports, ident)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.RBRACKET) {
		// TODO: Error handling
		return nil
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Successfully parsed module statement\n", p.encounterCount))
		p.log.AppendParser(fmt.Sprintf("\t- Module name: `%s`\n\t- ", stmt.Name.Value))
		p.log.AppendParser("Imports:\n")
		for _, imp := range stmt.Imports {
			p.log.AppendParser(fmt.Sprintf("\t\t- `%s`\n", imp.Value))
		}
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Steps in parsing let statement:\n", p.encounterCount))
		p.log.AppendParser(fmt.Sprintf("\n\ta. Assigning token to the statement to track positioning [line: %d, col: %d]\n", p.curToken.Line, p.curToken.Col))
	}
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		// TODO: Error handling
		return nil
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\tb. Assigning valid identifier `%s` to the let statement\n", p.curToken.Literal))
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		// TODO: Error handling
		return nil
	}

	p.nextToken()

	if p.debug {
		p.log.AppendParser("\n\tc. Parsing the expression to assign to the let statement...\n")
	}

	stmt.Value = p.parseExpression(LOWEST)

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\td. Successfully parsed a valid expression to assign to the let statement: `%s`\n", stmt.Value.String()))
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\te. Successfully parsed the entire let statement: `%s`\n", stmt.String()))
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Steps in parsing return statement:\n", p.encounterCount))
		p.log.AppendParser(fmt.Sprintf("\n\ta. Assigning token to the statement to track positioning [line: %d, col: %d]\n", p.curToken.Line, p.curToken.Col))
	}
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if p.debug {
		p.log.AppendParser("\n\tb. Parsing the expression to return...\n")
	}

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\tc. Successfully parsed a valid expression to return: `%s`\n", stmt.ReturnValue.String()))
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\td. Successfully parsed the entire return statement: `%s`\n", stmt.String()))
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		// TODO: Error handling
		return nil
	}

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		// TODO: Error handling
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	if !p.curTokenIs(token.RBRACE) {
		// TODO: Error handling
		return nil
	}

	p.nextToken()

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}
	
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	
	if !p.expectPeek(token.LET) {
		return nil
	}
	
	stmt.Init = p.parseLetStatement()
	
	if !p.curTokenIs(token.SEMICOLON) {
		return nil
	}
	p.nextToken()
	
	stmt.Condition = p.parseExpression(LOWEST)
	
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	p.nextToken()
	
	stmt.Post = p.parseExpression(LOWEST)
	fmt.Printf("post: %v\n", stmt.Post)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	if !p.curTokenIs(token.RBRACE) {
		return nil
	}

	fmt.Printf("stmt: %+v\n", stmt)

	return stmt
}

func (p *Parser) parseExpressionOrAssignStatement() ast.Statement {
	ident := p.parseIdentifier()
	if p.peekTokenIs(token.ASSIGN) {
		if p.debug {
			p.log.AppendParser(fmt.Sprintf("%d. Encountered an assignment statement, verifying whether the identifier is valid...\n", p.encounterCount))
		}
		if id, ok := ident.(*ast.Identifier); ok {
			if p.debug {
				p.log.AppendParser(fmt.Sprintf("\n\t- Identifier `%s` is valid, proceeding to parse the assignment statement\n", id.Value))
			}
			return p.parseAssignStatement(id)
		}
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Did not encounter an assign (`=`) token, so this is an expression statement\n", p.encounterCount))
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parseAssignStatement(ident *ast.Identifier) *ast.AssignStatement {
	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Steps in parsing assign statement:\n", p.encounterCount))
		p.log.AppendParser(fmt.Sprintf("\n\ta. Assigning token to the statement to track positioning [line: %d, col: %d]\n", p.curToken.Line, p.curToken.Col))
	}
	stmt := &ast.AssignStatement{Token: p.peekToken, Name: ident}

	p.nextToken()
	p.nextToken()

	if p.debug {
		p.log.AppendParser("\n\tb. Parsing the expression to assign to the identifier...\n")
	}
	stmt.Value = p.parseExpression(LOWEST)

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\tc. Successfully parsed a valid expression to assign to the identifier: `%s`\n", stmt.Value.String()))
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\td. Successfully parsed the entire assign statement: `%s`\n", stmt.String()))
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Steps in parsing expression statement:\n", p.encounterCount))
		p.log.AppendParser(fmt.Sprintf("\n\ta. Assigning token to the statement to track positioning [line: %d, col: %d]\n", p.curToken.Line, p.curToken.Col))
	}
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	if p.debug {
		p.log.AppendParser("\n\tb. Parsing expression statements is extremely simple. Just parse the expression and wrap it within a Statement implementation...\n")
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\tc. Successfully parsed the expression statement: `%s`\n", stmt.String()))
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\td. Successfully parsed the entire expression statement: `%s`\n", stmt.String()))
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if p.debug {
		p.log.AppendParser(fmt.Sprintf("%d. Steps in parsing block statement:\n", p.encounterCount))
		p.log.AppendParser(fmt.Sprintf("\n\ta. Assigning token to the statement to track positioning [line: %d, col: %d]\n", p.curToken.Line, p.curToken.Col))
	}
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	if p.debug {
		p.log.AppendParser("\n\tb. Parsing block statements is pretty simple. We only need to loop through all the statements within the block and store them until we reach the end of the block, signified `}`\n")
	}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			if p.debug {
				p.log.AppendParser(fmt.Sprintf("\n\t- Successfully parsed a statement to append to the block's `Statements` slice: `%s`\n", stmt.String()))
			}
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	if p.debug {
		p.log.AppendParser(fmt.Sprintf("\n\tc. Successfully parsed the entire block statement: `%s`\n", block.String()))
	}

	return block
}
