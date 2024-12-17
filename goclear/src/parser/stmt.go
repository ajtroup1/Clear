package parser

import (
	"github.com/ajtroup1/goclear/src/ast"
	"github.com/ajtroup1/goclear/src/lexer"
)

func parse_stmt (p *parser) ast.Statement {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expr := parse_expr(p, defalt_bp)
	p.expect(lexer.SEMI)

	return ast.ExpressionStatement{
		Expr: expr,
	}
}
