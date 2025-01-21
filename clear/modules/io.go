package modules

import (
	"fmt"

	"github.com/ajtroup1/clear/object"
)

var IOBuiltins = map[string]*object.Builtin{
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return &object.Null{}
		},
	},

	"println": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			out := ""
			for _, arg := range args {
				out += arg.Inspect()
			}
			fmt.Println(out)
			return &object.Null{}
		},
	},

	"printf": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.Error{Message: "printf requires at least one argument"}
			}

			format := args[0].(*object.String).Value

			if len(args) == 1 {
				fmt.Print(format)
				return &object.Null{}
			}

			values := make([]interface{}, len(args)-1)
			for i, arg := range args[1:] {
				switch arg := arg.(type) {
				case *object.Integer:
					values[i] = arg.Value
				case *object.Float:
					values[i] = arg.Value
				case *object.String:
					values[i] = arg.Value
				case *object.Boolean:
					values[i] = arg.Value
				default:
					values[i] = arg.Inspect()
				}
			}

			fmt.Printf(format, values...)
			return &object.Null{}
		},
	},

	"input": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			var input string
			fmt.Scanln(&input)
			return &object.String{Value: input}
		},
	},
}
