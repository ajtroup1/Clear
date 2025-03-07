package modules

import (
	"github.com/ajtroup1/clear/object"
)

var OSBuiltins = map[string]*object.Builtin{
	"exit": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.INTEGER_OBJ {
				return &object.Error{Message: "argument must be INTEGER"}
			}

			exitCode := args[0].(*object.Integer).Value

			return &object.ReturnValue{Value: &object.Integer{Value: exitCode}}
		},
	},
}
