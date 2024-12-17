package parser

import (
	"fmt"
	"strconv"

	"github.com/ajtroup1/goclear/src/ast"
	"github.com/ajtroup1/goclear/src/lexer"
	// "github.com/ajtroup1/goclear/src/util"
)

// * EXAMPLE PARSE: "10 + 5 * 2" *
// INITIAL parse_expr() iteration:
// Initially, parse_expr is called with a LOWEST binding power

// * SECOND ITERATION parse_expr() comments are highlighted (using VSCode's 'Better Comments' extension):
// * SECOND iteration comments are also prefixed with '*' as a result
// * Now called with a bp of ADDITIVE (5)
// * The parser is now currently investigating the 5 token since parse_binary_expr() consumed the '+'
// * AT THIS POINT, we have a binary expression:
// * ast.BinaryExpr {
// * 	Left: 10,
// * 	Operator: +,
// *	Right: ???, // This is what we are now assigning
// * }
// * So, it would be helpful to only think about this iteration as parsing "5 * 2" and assigning that to the right side of "10 +" later

// *       +
// *      / \
// *    10   *
// *        / \
// *       5   2

func parse_expr (p *parser, bp binding_power) ast.Expression {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s\n", lexer.TokenTypeString(tokenKind)))
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s\n", lexer.TokenTypeString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left
}

func parse_binary_expr (p *parser, left ast.Expression, bp binding_power) ast.Expression {
	operatorToken := p.advance()
	right := parse_expr(p, defalt_bp)

	return ast.BinaryExpr{
		Left: left,
		Operator: operatorToken,
		Right: right,
	}
}

func parse_primary_expr (p *parser) ast.Expression {
	switch	p.currentTokenKind() {
		case lexer.NUMBER:
			number, _ := strconv.ParseFloat(p.advance().Literal, 64)
			return ast.NumberExpr{
				Value: number,
			}
		case lexer.STRING:
			return ast.StringExpr{
				Value: p.advance().Literal,
			}
		case lexer.IDENT:
			return ast.SymbolExpr{
				Value: p.advance().Literal,
			}
		default:
			panic(fmt.Sprintf("Cannot create primary_expr from %s\n", lexer.TokenTypeString(p.currentTokenKind())))
	}
}
