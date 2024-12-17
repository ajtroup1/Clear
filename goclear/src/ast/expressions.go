package ast

import "github.com/ajtroup1/goclear/src/lexer"

// -------------------
// LITERAL EXPRESSIONS
// -------------------

type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

// -------------------
// COMPLEX EXPRESSIONS
// -------------------

type BinaryExpr struct {
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

func (n BinaryExpr) expr() {}
