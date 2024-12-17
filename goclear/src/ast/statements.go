package ast

type BlockStmt struct {
	Body []Statement
}

func (n BlockStmt) stmt() {}

type ExpressionStatement struct {
	Expr Expression
}

func (n ExpressionStatement) stmt() {}
