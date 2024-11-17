/*
	Defines the nodes necessary for constructing and storing data in the AST
	The AST will act as a tree structure for the evaluator to interpret recursively
	Each node acts differently, storing unique data necessary for its respective 'instruction'
	The 'root' of the AST is Program, which contains a list of statements
		Everything in Clear in encapsulated into statements somehow
		So they can be added to a slice of statements in the root Program node
*/

package ast

import (
	"bytes"

	"github.com/ajtroup1/clear/src/token"
)

// Base interface for every node in Clear
// Every node must be able to return a string of its literal value and a string of its representation
type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Structure for a Clear return statement
// Simply stores a single Expression, later evaluated to a value, to return from a function or whatever
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// Wraps an expression within a statement wrapper in order to envelope it in Program's Statements slice
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// Prefix expressions are simple ome-sided expressions such as
	// Not (!), Negate (-), ....
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// Infix expressions contain a left and right value with their operand, unlike prefix expressions
	// 1 + 2, 3 * (7 - 2), 8 + 8 + 8, ....
// Infix Expressions can be "infinitely" large since since it encapsulates 2 expressions, not values
	// Since "5" is an expression, you can have "5 + 5"
	// Since "5 + 5" is also an expression, you can also have "5 + 5 * 2"
		// Here, "5 * 2" would be the right side of the expression (due to precedence), and "5" is the left side
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}
