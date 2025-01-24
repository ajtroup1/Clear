package ast

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/compiled_clear/src/token"
)

type Node interface {
	TokenLiteral() string
	ToString() string
	Position() (int, int)
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Node
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) ToString() string {
	var buf strings.Builder
	for _, s := range p.Statements {
		buf.WriteString(s.ToString())
	}
	return buf.String()
}
func (p *Program) Position() (int, int) {
	return 0, 0
}

// ----------------------------------------------------------------------------
// STATEMENTS
// ----------------------------------------------------------------------------

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) ToString() string {
	var buf strings.Builder
	buf.WriteString(ls.TokenLiteral() + " " + ls.Name.Value + " = ")
	buf.WriteString(ls.Value.ToString())
	buf.WriteString(";")
	return buf.String()
}
func (ls *LetStatement) Position() (int, int) {
	return ls.Token.Line, ls.Token.Col
}

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) ToString() string {
	var buf strings.Builder
	buf.WriteString(rs.TokenLiteral() + " ")
	buf.WriteString(rs.ReturnValue.ToString())
	buf.WriteString(";")
	return buf.String()
}
func (rs *ReturnStatement) Position() (int, int) {
	return rs.Token.Line, rs.Token.Col
}

// ----------------------------------------------------------------------------
// EXPRESSIONS
// ----------------------------------------------------------------------------

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) ToString() string {
	var buf strings.Builder
	buf.WriteString("(")
	buf.WriteString(pe.Operator)
	buf.WriteString(pe.Right.ToString())
	buf.WriteString(")")
	buf.WriteString(";")
	return buf.String()
}
func (pe *PrefixExpression) Position() (int, int) {
	return pe.Token.Line, pe.Token.Col
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) ToString() string {
	var buf strings.Builder
	buf.WriteString("(")
	buf.WriteString(ie.Left.ToString())
	buf.WriteString(" " + ie.Operator + " ")
	buf.WriteString(ie.Right.ToString())
	buf.WriteString(")")
	buf.WriteString(";")
	return buf.String()
}
func (ie *InfixExpression) Position() (int, int) {
	return ie.Token.Line, ie.Token.Col
}

type PostfixExpression struct {
	Token    token.Token // The postfix token, e.g. ++
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) ToString() string {
	var buf strings.Builder
	buf.WriteString("(")
	buf.WriteString(pe.Left.ToString())
	buf.WriteString(pe.Operator)
	buf.WriteString(")")
	buf.WriteString(";")
	return buf.String()
}
func (pe *PostfixExpression) Position() (int, int) {
	return pe.Token.Line, pe.Token.Col
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) ToString() string     { return i.Value }
func (i *Identifier) Position() (int, int) {
	return i.Token.Line, i.Token.Col
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) ToString() string {
	return fmt.Sprintf("%d", il.Value)
}
func (il *IntegerLiteral) Position() (int, int) {
	return il.Token.Line, il.Token.Col
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) ToString() string {
	return fmt.Sprintf("%f", fl.Value)
}
func (fl *FloatLiteral) Position() (int, int) {
	return fl.Token.Line, fl.Token.Col
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) ToString() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) Position() (int, int) {
	return b.Token.Line, b.Token.Col
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) ToString() string {
	return sl.Value
}
func (sl *StringLiteral) Position() (int, int) {
	return sl.Token.Line, sl.Token.Col
}
