package ast

import (
	"fmt"

	"github.com/ajtroup1/goclear/token"
)

type Node interface {
	ToString() string
	Position() (line, col int)
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

type Program struct {
	Statements []Statement
}

func (p *Program) ToString() string {
	str := ""
	for _, stmt := range p.Statements {
		str += stmt.ToString() + "\n"
	}
	return str
}
func (p *Program) Position() (line, col int) {
	return 0, 0
}

type BaseNode struct {
	Token token.Token
}

func (bn *BaseNode) Position() (line, col int) {
	return bn.Token.Line, bn.Token.Col
}

// ==========
// STATEMENTS
// ==========

type BlockStatement struct {
	BaseNode
	Statements []Statement
}

func (bs *BlockStatement) statement() {}
func (bs *BlockStatement) ToString() string {
	str := ""
	for _, stmt := range bs.Statements {
		str += stmt.ToString() + "\n"
	}
	return str
}

type LetStatement struct {
	BaseNode
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statement() {}
func (ls *LetStatement) ToString() string {
	return fmt.Sprintf("LET %s = %v", ls.Name.Value, ls.Value)
}

type ReturnStatement struct {
	BaseNode
	Value Expression
}

func (rs *ReturnStatement) statement() {}
func (rs *ReturnStatement) ToString() string {
	return fmt.Sprintf("RETURN %v", rs.Value)
}

type ExpressionStatement struct {
	BaseNode
	Expression Expression
}

func (es *ExpressionStatement) statement() {}
func (es *ExpressionStatement) ToString() string {
	return es.Expression.ToString()
}

// ===========
// EXPRESSIONS
// ===========

type Identifier struct {
	BaseNode
	Value  string
}

func (i *Identifier) expression() {}
func (i *Identifier) ToString() string {
	return fmt.Sprintf("IDENT %s", i.Value)
}

type IntegerLiteral struct {
	BaseNode
	Value int64
}

func (i *IntegerLiteral) expression() {}
func (i *IntegerLiteral) ToString() string {
	return fmt.Sprintf("INT %d", i.Value)
}

type FloatLiteral struct {
	BaseNode
	Value float64
}

func (f *FloatLiteral) expression() {}
func (f *FloatLiteral) ToString() string {
	return fmt.Sprintf("FLOAT %f", f.Value)
}

type StringLiteral struct {
	BaseNode
	Value string
}

func (s *StringLiteral) expression() {}
func (s *StringLiteral) ToString() string {
	return fmt.Sprintf("STRING %s", s.Value)
}

type PrefixExpression struct {
	BaseNode
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expression() {}
func (pe *PrefixExpression) ToString() string {
	return fmt.Sprintf("(%s) %v", pe.Operator, pe.Right)
}

type InfixExpression struct {
	BaseNode
	Operator string
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) expression() {}
func (ie *InfixExpression) ToString() string {
	return fmt.Sprintf("(%v) %s (%v)", ie.Left, ie.Operator, ie.Right)
}

type Boolean struct {
	BaseNode
	Value bool
}

func (b *Boolean) expression() {}
func (b *Boolean) ToString() string {
	return fmt.Sprintf("BOOL %t", b.Value)
}

type IfExpression struct {
	BaseNode
	Condition       Expression
	Consequence     *BlockStatement
	Alternative     *BlockStatement
}

func (ie *IfExpression) expression() {}
func (ie *IfExpression) ToString() string {
	str := fmt.Sprintf("IF %v %v", ie.Condition, ie.Consequence)
	if ie.Alternative != nil {
		str += fmt.Sprintf(" ELSE %v", ie.Alternative)
	}
	return str
}

type FunctionLiteral struct {
	BaseNode
	Name 		 	*Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expression() {}
func (fl *FunctionLiteral) ToString() string {
	return fmt.Sprintf("FUNCTION %v %v", fl.Parameters, fl.Body)
}

type CallExpression struct {
	BaseNode
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expression() {}
func (ce *CallExpression) ToString() string {
	return fmt.Sprintf("CALL %v %v", ce.Function, ce.Arguments)
}