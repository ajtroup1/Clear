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

	"input": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			var input string
			fmt.Scanln(&input)
			return &object.String{Value: input}
		},
	},
}
