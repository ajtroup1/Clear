package ast

import (
	"bytes"

	"github.com/ajtroup1/goclear/src/token"
)

// Generic interface for every node in the AST
// Every node in the AST must be able to:
// - `TokenLiteral()` Return the literal value associated with the node
// - `String()` Return a string representation of the node
type Node interface {
	TokenLiteral() string
	String() string
}

// Generic interface for Statements in Clear
// Uses marker method statementNode()
type Statement interface {
	Node
	statementNode()
}

// Generic interface for Expressions in Clear
// Uses marker method expressionNode()
type Expression interface {
	Node
	expressionNode()
}

// Root node of the AST, containing every node in the program within
// Everything in the program can be enveloped within a statement, so a slice of statements
// forms the Program node
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	} else {
		return p.Statements[0].TokenLiteral()
	}
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Statement to assign variables in Clear, uses `let`
// let x = 7;
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	// The value being assigned could be anything, which is why Expression is being used
	// ex. '5', '5 + 2', 'myFunction(param)', 'myObject.property'
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

// Name identifier for variables
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Return statement which exits from the current scope
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

// It is necessary to encapsulate expressions into statements to include them into the Statements[] heirarchy
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

// Simple integer node to store an int64
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// Prefix expressions have an operator to the left of the expression
// There is no need for a "left hand side" in prefix expressions
type PrefixExpression struct {
	Token    token.Token // -, !, ...
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

// Most mathematical expressions are infix
// 	"1 + 2"
// 	"2 * func() - 5"
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

// Simple boolean structure
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// Node for if EXPRESSIONS in Clear
// These if expressions return a value (true, false), meaning they are an expression
// If the "if" was only used for control flow (Go, C, ...), then it would be a statement
// This is also really an "if / else" expression since it provides an alternative BlockStatement
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression  // if --> (x > 5) <-- Condition
	Consequence *BlockStatement // if (x > 5) { do this < -- Consequence }
	Alternative *BlockStatement // if (x > 5) { do this } else { do that <-- Alternative }
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// Simply a wrapper around a slice of statements
// Necessary for if consequences/alternatives and functions
// Program is also really a Block Statement
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
