package modules

import (
	"fmt"
	"math"

	"github.com/ajtroup1/clear/object"
)

var MathBuiltins = map[string]*object.Builtin{
	"abs": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}

			if args[0].Type() == object.INTEGER_OBJ {
				return &object.Integer{Value: int64(math.Abs(float64(args[0].(*object.Integer).Value)))}
			}

			if args[0].Type() == object.FLOAT_OBJ {
				return &object.Float{Value: math.Abs(args[0].(*object.Float).Value)}
			}

			fmt.Printf("args[0].Type(): %s\n", args[0].Type())

			return &object.Error{Message: fmt.Sprintf("argument to `abs` not supported, got %s", args[0].Type())}
		},
	},

	"round": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}

			if args[0].Type() == object.INTEGER_OBJ {
				return args[0]
			}

			if args[0].Type() == object.FLOAT_OBJ {
				return &object.Integer{Value: int64(math.Round(args[0].(*object.Float).Value))}
			}

			return &object.Error{Message: fmt.Sprintf("argument to `round` not supported, got %s", args[0].Type())}
		},
	},

	"pow": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}

			if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
				return &object.Integer{Value: int64(math.Pow(float64(args[0].(*object.Integer).Value), float64(args[1].(*object.Integer).Value)))}
			}

			if args[0].Type() == object.FLOAT_OBJ && args[1].Type() == object.FLOAT_OBJ {
				return &object.Float{Value: math.Pow(args[0].(*object.Float).Value, args[1].(*object.Float).Value)}
			}

			if args[0].Type() == object.FLOAT_OBJ && args[1].Type() == object.INTEGER_OBJ {
				return &object.Float{Value: math.Pow(args[0].(*object.Float).Value, float64(args[1].(*object.Integer).Value))}
			}

			return &object.Error{Message: fmt.Sprintf("arguments to `pow` not supported, got %s and %s", args[0].Type(), args[1].Type())}
		},
	},
}
