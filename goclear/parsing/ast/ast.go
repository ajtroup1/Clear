package ast

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/goclear/lexing/token"
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
	Imports    []*ModuleStatement
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

type AssignStatement struct {
	BaseNode
	Name  *Identifier
	Value Expression
	Type  token.TokenType
}

func (as *AssignStatement) statement() {}
func (as *AssignStatement) ToString() string {
	return fmt.Sprintf("LET %s = %v", as.Name.Value, as.Value)
}

type ConstStatement struct {
	BaseNode
	Name  *Identifier
	Value Expression
	Type  token.TokenType
}

func (cs *ConstStatement) statement() {}
func (cs *ConstStatement) ToString() string {
	return fmt.Sprintf("CONST %s = %v", cs.Name.Value, cs.Value)
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

type WhileStatement struct {
	BaseNode
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statement() {}
func (ws *WhileStatement) ToString() string {
	return fmt.Sprintf("WHILE %v %v", ws.Condition, ws.Body)
}

type ForStatement struct {
	BaseNode
	Init      Statement
	Condition Expression
	Post      Expression
	Body      *BlockStatement
}

func (fs *ForStatement) statement() {}
func (fs *ForStatement) ToString() string {
	return fmt.Sprintf("FOR %v %v %v %v", fs.Init, fs.Condition, fs.Post, fs.Body)
}

type BreakStatement struct {
	BaseNode
}

func (bs *BreakStatement) statement() {}
func (bs *BreakStatement) ToString() string {
	return "BREAK"
}

type ContinueStatement struct {
	BaseNode
}

func (cs *ContinueStatement) statement() {}
func (cs *ContinueStatement) ToString() string {
	return "CONTINUE"
}

type ModuleStatement struct {
	BaseNode
	Name string
}

func (is *ModuleStatement) statement() {}
func (is *ModuleStatement) ToString() string {
	return fmt.Sprintf("IMPORT %s", is.Name)
}

// ===========
// EXPRESSIONS
// ===========

type Identifier struct {
	BaseNode
	Value string
	Type  token.TokenType
}

func (i *Identifier) expression() {}
func (i *Identifier) ToString() string {
	return fmt.Sprintf("IDENT %s: %s", i.Value, strings.ToLower(string(i.Type)))
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

type CharLiteral struct {
	BaseNode
	Value rune
}

type BooleanLiteral struct {
	BaseNode
	Value bool
}

func (b *BooleanLiteral) expression() {}
func (b *BooleanLiteral) ToString() string {
	return fmt.Sprintf("BOOL %t", b.Value)
}

func (c *CharLiteral) expression() {}
func (c *CharLiteral) ToString() string {
	return fmt.Sprintf("CHAR %c", c.Value)
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

type PostfixExpression struct {
	BaseNode
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expression() {}
func (pe *PostfixExpression) ToString() string {
	return fmt.Sprintf("%v (%s)", pe.Left, pe.Operator)
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

type IfExpression struct {
	BaseNode
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
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
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	ReturnType token.TokenType
}

func (fl *FunctionLiteral) expression() {}
func (fl *FunctionLiteral) ToString() string {
	return fmt.Sprintf("FUNCTION:\n\tReturn Type: %v\n\tParameters: %v\n\tBody: %v", fl.ReturnType, fl.Parameters, fl.Body)
}

type CallExpression struct {
	BaseNode
	FunctionIdentifier  Expression
	Arguments []CallArgument
}

func (ce *CallExpression) expression() {}
func (ce *CallExpression) ToString() string {
	return fmt.Sprintf("CALL %s\nArguments:\n\t%v", ce.FunctionIdentifier, ce.Arguments)
}

type CallArgument struct {
	BaseNode
	Expression Expression
	Type       token.TokenType
}

func (ca *CallArgument) expression() {}
func (ca *CallArgument) ToString() string {
	return fmt.Sprintf("ARG %v %s", ca.Expression, ca.Type)
}

type ClassStatement struct {
	BaseNode
	Name       *Identifier
	Properties []*PropertyStatement
	Methods    []*FunctionLiteral
}

func (cs *ClassStatement) statement() {}
func (cs *ClassStatement) ToString() string {
	return fmt.Sprintf("CLASS %s %v %v", cs.Name.Value, cs.Properties, cs.Methods)
}

type PropertyStatement struct {
	BaseNode
	Name  string
	Type  string
	Value interface{}
}

func (ps *PropertyStatement) statement() {}
func (ps *PropertyStatement) ToString() string {
	return fmt.Sprintf("PROPERTY %s %s %v", ps.Name, ps.Type, ps.Value)
}
