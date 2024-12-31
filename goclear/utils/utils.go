package utils

import (
	"github.com/ajtroup1/goclear/parsing/ast"
	"github.com/sanity-io/litter"
)

func PrettyPrintASTNode(node ast.Node) {
	litter.Dump(node)
}