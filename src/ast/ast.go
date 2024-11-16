/*
	Defines the nodes necessary for constructing and storing data in the AST
	The AST will act as a tree structure for the evaluator to interpret recursively
	Each node acts differently, storing unique data necessary for its respective 'instruction'
	The 'root' of the AST is Program, which contains a list of statements
		Everything in Clear in encapsulated into statements somehow
		So they can be added to a slice of statements in the root Program node
*/

package ast

import "github.com/ajtroup1/clear/src/token"

// Base interface for every node in Clear
// Every node must be able to return a string of its literal value
type Node interface {
	TokenLiteral() string
}
// Base interface for every statement in Clear
// Every statement is a node and implements the tag statementNode()
type Statement interface {
	Node
	statementNode()
}
// Base interface for every expression in Clear
// Every statement is a node and implements the tag expressionNode()
type Expression interface {
	Node
	expressionNode()
}

// Root node of the AST
// Contains a list of statements that encapsulates all the code in the program
type Program struct {
	Statements []Statement
}

// TokenLiteral() implementation for Program
// Returns the first statement of the Program, giving a rough idea of the program's starting point
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Structure for a Clear let statement, or variable assignment
// Stores a name for the identifier to go by and a value to set it equal to
	// Value is an Expression because every value is an expression: "5 + 2", "(5+2) * 7", "5" <-- Even integer literals
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// Structure for a Clear return statement
// Simply stores a single Expression, later evaluated to a value, to return from a function or whatever
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
