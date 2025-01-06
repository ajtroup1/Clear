package ast

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/goclear/lexing/token"
)

type DataType string

const (
	UNKNOWN  DataType = "UNKNOWN"
	INT      DataType = "INT"
	FLOAT    DataType = "FLOAT"
	STRING   DataType = "STRING"
	CHAR     DataType = "CHAR"
	BOOL     DataType = "BOOL"
	VOID     DataType = "VOID"
	FUNCTION DataType = "FUNCTION"
	MODULE   DataType = "MODULE"
	// MODULEFUNCTION is a special type used to represent a function imported from a module
	// This is used to differentiate between functions defined in the current file and functions imported from a module
	// Example: "mod math: Round" would be of type MODULEFUNCTION and it's symbol would look like:
	// {
	// 	"Name": "Round",
	// 	"Type": "MODULEFUNCTION",
	// 	"Value": "math" <-- The module name. Uses the Value property especially for the parent import identifier
	// }
	MODULEFUNCTION DataType = "MODULEFUNCTION"
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
	GetType() DataType
}

type Program struct {
	Statements []Statement        `json:"Statements"`
	Imports    []*ModuleStatement `json:"Imports"`
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
	Statements []Statement `json:"Statements"`
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
	Name  *Identifier `json:"Name"`
	Value Expression  `json:"Value"`
	Type  DataType    `json:"Type"`
}

func (as *AssignStatement) statement() {}
func (as *AssignStatement) ToString() string {
	return fmt.Sprintf("(%s) %s = %v", as.Type, as.Name.Value, as.Value.ToString())
}

type ReassignStatement struct {
	BaseNode
	Name  *Identifier `json:"Name"`
	Value Expression  `json:"Value"`
}

func (rs *ReassignStatement) statement() {}
func (rs *ReassignStatement) ToString() string {
	return fmt.Sprintf("REASSIGN (%s) %s = %v", rs.Name.Type, rs.Name.Value, rs.Value.ToString())
}

type ConstStatement struct {
	BaseNode
	Name  *Identifier `json:"Name"`
	Value Expression  `json:"Value"`
	Type  DataType    `json:"Type"`
}

func (cs *ConstStatement) statement() {}
func (cs *ConstStatement) ToString() string {
	return fmt.Sprintf("CONST %s = %v", cs.Name.Value, cs.Value)
}

type ReturnStatement struct {
	BaseNode
	Value Expression `json:"Value"` // The value to return
}

func (rs *ReturnStatement) statement() {}
func (rs *ReturnStatement) ToString() string {
	return fmt.Sprintf("RETURN %v", rs.Value.ToString())
}

type ExpressionStatement struct {
	BaseNode
	Expression Expression `json:"Expression"`
}

func (es *ExpressionStatement) statement() {}
func (es *ExpressionStatement) ToString() string {
	return es.Expression.ToString()
}

type WhileStatement struct {
	BaseNode
	Condition Expression      `json:"Condition"`
	Body      *BlockStatement `json:"Body"`
}

func (ws *WhileStatement) statement() {}
func (ws *WhileStatement) ToString() string {
	return fmt.Sprintf("WHILE %v %v", ws.Condition, ws.Body)
}

type ForStatement struct {
	BaseNode
	Init      Statement       `json:"Init"`
	Condition Expression      `json:"Condition"`
	Post      Expression      `json:"Post"`
	Body      *BlockStatement `json:"Body"`
}

func (fs *ForStatement) statement() {}
func (fs *ForStatement) ToString() string {
	return fmt.Sprintf("FOR:\n\tInit: %v\n\tCondition: %v\n\tPost: %v\n\tBody: %v", fs.Init.ToString(), fs.Condition.ToString(), fs.Post.ToString(), fs.Body.ToString())
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
	Name *Identifier `json:"Name"` // The name of the module
	// Flag to determine if all functions should be imported
	// Toggled to true by using '*': "mod math: *"
	ImportAll bool `json:"ImportAll"`
	// If not importing all, this will be a list of functions to import
	// Contained within brackets and separated by commas: "mod math: [Round, Pow]"
	Imports []*Identifier `json:"Imports"`
}

func (is *ModuleStatement) statement() {}
func (is *ModuleStatement) ToString() string {
	output := fmt.Sprintf("IMPORT %s\n\tImport all: %v\n\tImports: ", is.Name.ToString(), is.ImportAll)
	output += "[\n\t\t"
	for _, imp := range is.Imports {
		output += "   " + imp.ToString() + ",\n\t\t"
	}
	output += " ]"
	return output
}

type ClassStatement struct {
	BaseNode
	Name       *Identifier          `json:"Name"`
	Properties []*PropertyStatement `json:"Properties"`
	Methods    []*FunctionLiteral   `json:"Methods"`
}

func (cs *ClassStatement) statement() {}
func (cs *ClassStatement) ToString() string {
	return fmt.Sprintf("CLASS %s %v %v", cs.Name.Value, cs.Properties, cs.Methods)
}

type PropertyStatement struct {
	BaseNode
	Name  string      `json:"Name"`
	Type  DataType    `json:"Type"`
	Value interface{} `json:"Value"`
}

func (ps *PropertyStatement) statement() {}
func (ps *PropertyStatement) ToString() string {
	return fmt.Sprintf("PROPERTY %s %s %v", ps.Name, ps.Type, ps.Value)
}

// ===========
// EXPRESSIONS
// =========== 

type Identifier struct {
	BaseNode
	Value string   `json:"value"`
	Type  DataType `json:"type"`
}

func (i *Identifier) expression() {}
func (i *Identifier) ToString() string {
	return fmt.Sprintf("IDENT %s: (%s)", i.Value, strings.ToLower(string(i.Type)))
}
func (i *Identifier) GetType() DataType {
	return i.Type
}

type IntegerLiteral struct {
	BaseNode
	Value int64 `json:"value"`
}

func (i *IntegerLiteral) expression() {}
func (i *IntegerLiteral) ToString() string {
	return fmt.Sprintf("INT %d", i.Value)
}
func (i *IntegerLiteral) GetType() DataType {
	return INT
}

type FloatLiteral struct {
	BaseNode
	Value float64 `json:"value"`
}

func (f *FloatLiteral) expression() {}
func (f *FloatLiteral) ToString() string {
	return fmt.Sprintf("FLOAT %f", f.Value)
}
func (f *FloatLiteral) GetType() DataType {
	return FLOAT
}

type StringLiteral struct {
	BaseNode
	Value string `json:"value"`
}

func (s *StringLiteral) expression() {}
func (s *StringLiteral) ToString() string {
	return fmt.Sprintf("STRING %s", s.Value)
}
func (s *StringLiteral) GetType() DataType {
	return STRING
}

type CharLiteral struct {
	BaseNode
	Value rune `json:"value"`
}

func (c *CharLiteral) expression() {}
func (c *CharLiteral) ToString() string {
	return fmt.Sprintf("CHAR %c", c.Value)
}
func (c *CharLiteral) GetType() DataType {
	return CHAR
}

type BooleanLiteral struct {
	BaseNode
	Value bool `json:"value"`
}

func (b *BooleanLiteral) expression() {}
func (b *BooleanLiteral) ToString() string {
	return fmt.Sprintf("BOOL %t", b.Value)
}
func (b *BooleanLiteral) GetType() DataType {
	return BOOL
}

type PrefixExpression struct {
	BaseNode
	Operator string     `json:"operator"`
	Right    Expression `json:"right"`
}

func (pe *PrefixExpression) expression() {}
func (pe *PrefixExpression) ToString() string {
	return fmt.Sprintf("(%s) %v", pe.Operator, pe.Right.ToString())
}
func (pe *PrefixExpression) GetType() DataType {
	if pe.Operator == "!" {
		return BOOL
	} else if pe.Operator == "-" {
		return pe.Right.GetType()
	}
	return UNKNOWN
}

type PostfixExpression struct {
	BaseNode
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
}

func (pe *PostfixExpression) expression() {}
func (pe *PostfixExpression) ToString() string {
	return fmt.Sprintf("%v (%s)", pe.Left.ToString(), pe.Operator)
}
func (pe *PostfixExpression) GetType() DataType {
	return UNKNOWN
}

type InfixExpression struct {
	BaseNode
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
	Right    Expression `json:"right"`
}

func (ie *InfixExpression) expression() {}
func (ie *InfixExpression) ToString() string {
	return fmt.Sprintf("(%v) %s (%v)", ie.Left.ToString(), ie.Operator, ie.Right.ToString())
}
func (ie *InfixExpression) GetType() DataType {
	lType := ie.Left.GetType()
	rType := ie.Right.GetType()
	if lType == rType {
		return lType
	}
	return UNKNOWN
}

type IfExpression struct {
	BaseNode
	Condition   Expression      `json:"condition"`
	Consequence *BlockStatement `json:"consequence"`
	Alternative *BlockStatement `json:"alternative,omitempty"`
}

func (ie *IfExpression) expression() {}
func (ie *IfExpression) ToString() string {
	str := fmt.Sprintf("IF %v %v", ie.Condition, ie.Consequence)
	if ie.Alternative != nil {
		str += fmt.Sprintf(" ELSE %v", ie.Alternative)
	}
	return str
}
func (ie *IfExpression) GetType() DataType {
	return BOOL
}

type FunctionLiteral struct {
	BaseNode
	Name       *Identifier     `json:"name"`
	Parameters []*Identifier   `json:"parameters"`
	Body       *BlockStatement `json:"body"`
	ReturnType DataType        `json:"returnType"`
}

func (fl *FunctionLiteral) expression() {}
func (fl *FunctionLiteral) ToString() string {
	paramString := ""
	for _, param := range fl.Parameters {
		paramString += param.ToString() + " "
	}
	return fmt.Sprintf("FUNCTION:\n\tReturn Type: %v\n\tParameters: %v\n\tBody: %v", fl.ReturnType, paramString, fl.Body.ToString())
}
func (fl *FunctionLiteral) GetType() DataType {
	return fl.ReturnType
}

type CallExpression struct {
	BaseNode
	Function  *FunctionLiteral `json:"function"`
	Arguments []CallArgument   `json:"arguments"`
}

func (ce *CallExpression) expression() {}
func (ce *CallExpression) ToString() string {
	return fmt.Sprintf("CALL %s\nArguments:\n\t%v", ce.Function.Name.Value, ce.Arguments)
}
func (ce *CallExpression) GetType() DataType {
	return ce.Function.GetType()
}

type CallArgument struct {
	BaseNode
	Expression Expression `json:"expression"`
	Type       DataType   `json:"type"`
}

func (ca *CallArgument) expression() {}
func (ca *CallArgument) ToString() string {
	return fmt.Sprintf("ARG %v %s", ca.Expression, ca.Type)
}
func (ca *CallArgument) GetType() DataType {
	return ca.Type
}

type TypeCastExpression struct {
	BaseNode
	Type  DataType   `json:"type"`  // The target type (e.g., "int", "float").
	Value Expression `json:"value"` // The expression being cast.
}

func (t *TypeCastExpression) expression()          {}
func (t *TypeCastExpression) TokenLiteral() string { return t.Token.Literal }
func (t *TypeCastExpression) ToString() string {
	return fmt.Sprintf("%s(%s)", t.Type, t.Value.ToString())
}
func (t *TypeCastExpression) GetType() DataType {
	return t.Type
}
