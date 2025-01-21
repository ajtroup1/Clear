/*
	The AST package simply contains all possible nodes in the Clear language
	as well as the structure and implementations of Statements & Expressions

	Nodes form into a structured tree of program information that can be
	traversed and evaluated by the interpreter

	Nodes are divided into Statements and Expressions, which is determined by
	whether the node returns a value or not
		- Statements do not return values, they dictate control flow
		- Expressions are evaluated to produce a result that be assigned or interpreted how you please
*/

package ast

import (
	"bytes"
	"strings"

	"github.com/ajtroup1/clear/token"
)

// The base Node interface
// Every node must be able to:
//   - Return its token's literal value
//   - Return a string representation of itself (aka ToString())
type Node interface {
	TokenLiteral() string
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode() // Tracking method for statements
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode() // Tracking method for expressions
}

// The root node of the AST, Program
// Any program simply consists of a series of statements,
// no matter how complex those statements become
type Program struct {
	NoStatements bool

	Statements []Statement        `json:"statements"`
	Modules    []*ModuleStatement `json:"modules"`
}

// Just return the literal value of the first statement
// Not very helpful, but TokenLiteral() is silly to call on Program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Since all statements implement the String() method, just call
// that on all Statements in the Program node
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// -------------------------
// # STATEMENT NODES

type ModuleStatement struct {
	Token     token.Token   // the 'module' token
	Name      *Identifier   `json:"name"`
	ImportAll bool          `json:"import_all"`
	Imports   []*Identifier `json:"imports"`
}

func (ms *ModuleStatement) statementNode()       {}
func (ms *ModuleStatement) TokenLiteral() string { return ms.Token.Literal }
func (ms *ModuleStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ms.TokenLiteral() + " ")
	out.WriteString(ms.Name.String())
	out.WriteString(" ")

	if ms.ImportAll {
		out.WriteString("*")
	} else {
		imports := []string{}
		for _, i := range ms.Imports {
			imports = append(imports, i.String())
		}
		out.WriteString(strings.Join(imports, ", "))
	}

	out.WriteString(";")

	return out.String()
}

// Statement used to assign variables to an expression
// let x = 7;
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier `json:"name"`
	Value Expression  `json:"value"`
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

type AssignStatement struct {
	Token token.Token
	Name  *Identifier `json:"name"`
	Value Expression  `json:"value"`
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(" = ")
	out.WriteString(as.Value.String())
	out.WriteString(";")

	return out.String()
}

// Statement used to return a value from a function
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression  `json:"returnValue"`
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

type WhileStatement struct {
	Token     token.Token
	Condition Expression      `json:"condition"`
	Body      *BlockStatement `json:"body"`
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ws.TokenLiteral() + " ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())

	return out.String()
}

type ForStatement struct {
	Token     token.Token
	Init      Statement       `json:"init"`
	Condition Expression      `json:"condition"`
	Post      Expression      `json:"post"`
	Body      *BlockStatement `json:"body"`
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fs.TokenLiteral() + " (")
	out.WriteString(fs.Init.String() + " ")
	out.WriteString(fs.Condition.String() + "; ")
	out.WriteString(fs.Post.String())
	out.WriteString(") {")
	out.WriteString(fs.Body.String())
	out.WriteString("}")

	return out.String()
}

// Since lone expressions cannot be enveloped in the Program's
// Statements slice, we must wrap them in an ExpressionStatement
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression  `json:"expression"`
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Simply wrap a slice of statements in a block
// Useful for if-else statements, loops, functions, ...
// { <--- Usually indicates the start of a block of statements
//
//		let x = 7;
//		let y = 8;
//		return x + y;
//	}
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement `json:"statements"`
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

// -------------------------
// # EXPRESSION NODES

// Basic identifier node for variables, functions, ...
// Only needs to hold a string indiciating the string value of the ident name
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string      `json:"value"`
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Boolean struct {
	Token token.Token
	Value bool `json:"value"`
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64 `json:"value"`
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64 `json:"value"`
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string `json:"value"`
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Value }

// Prefix expressions contain a prefix operator and a right expression
// Ex. !true, -5, ...
// [-], [!], ... <-- actual operators
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string      `json:"operator"`
	Right    Expression  `json:"right"`
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

// Infix expressions contain a left expression, an operator, and a right expression
// This is more like a 'traditional' expression or equation you would think of
// Ex. 5 + 5, 5 * (5 + 5), ...
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression  `json:"left"`
	Operator string      `json:"operator"`
	Right    Expression  `json:"right"`
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

type PostfixExpression struct {
	Token    token.Token
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")

	return out.String()
}

// If expressions contain a condition, a consequence, and an alternative
// These are expressions since if returns a true or false value
// The result of the condition statement (true or false) determines whether the consequence or alternative is executed
type IfExpression struct {
	Token       token.Token     // The 'if' token
	Condition   Expression      `json:"condition"`
	Consequence *BlockStatement `json:"consequence"`
	Alternative *BlockStatement `json:"alternative"`
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

// Function literals contain a list of parameters and a block statement to execute
// Obivously, functions may have infinite parameters and statements within the block
type FunctionLiteral struct {
	Token      token.Token     // The 'fn' token
	Parameters []*Identifier   `json:"parameters"`
	Body       *BlockStatement `json:"body"`
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// Expression that invokes a predefined function with an optional list of arguments
type CallExpression struct {
	Token     token.Token  // The '(' token
	Function  Expression   `json:"function"` // can be an Identifier or a FunctionLiteral
	Arguments []Expression `json:"arguments"`
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token  // The '[' token
	Elements []Expression `json:"elements"`
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression  `json:"left"`  // The expression to the left of the index
	Index Expression  `json:"index"` // The index expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token // The '{' token
	Pairs map[Expression]Expression `json:"pairs"`
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// -------------------------
