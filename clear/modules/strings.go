package modules

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/clear/object"
)

var StringsBuiltins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got type %s", arg.Type())}
			}
		},
	},

	"concat": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want >= 2", len(args))}
			}

			for _, arg := range args {
				if arg.Type() != object.STRING_OBJ {
					return &object.Error{
						Message: fmt.Sprintf("arguments to `concat` must be STRING, received a %s from '%v'", arg.Type(), arg.Inspect()),
					}
				}
			}

			var output string
			for _, arg := range args {
				strArg := arg.(*object.String)
				output += strArg.Value
			}

			return &object.String{Value: output}
		},
	},

	"concatDelim": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want >= 2", len(args))}
			}

			delimArg := args[0]
			if delimArg.Type() != object.STRING_OBJ {
				return &object.Error{
					Message: fmt.Sprintf("delimiter must be a STRING, received a %s from '%v'", delimArg.Type(), delimArg.Inspect()),
				}
			}
			delimiter := delimArg.(*object.String).Value

			for _, arg := range args[1:] {
				if arg.Type() != object.STRING_OBJ {
					return &object.Error{
						Message: fmt.Sprintf("arguments to `concatDelim` must be STRING, received a %s from '%v'", arg.Type(), arg.Inspect()),
					}
				}
			}

			var output string
			for i, arg := range args[1:] {
				strArg := arg.(*object.String) 
				output += strArg.Value
				if i < len(args[1:])-1 { 
					output += delimiter
				}
			}

			return &object.String{Value: output}
		},
	},

	"split": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: fmt.Sprintf("first argument must be STRING, got %s", args[0].Type())}
			}

			if args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: fmt.Sprintf("second argument must be STRING, got %s", args[1].Type())}
			}

			strArg := args[0].(*object.String)
			delimiterArg := args[1].(*object.String)

			parts := strings.Split(strArg.Value, delimiterArg.Value)

			array := &object.Array{}
			for _, part := range parts {
				array.Elements = append(array.Elements, &object.String{Value: part})
			}

			return array
		},
	},

	"lower": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: fmt.Sprintf("argument must be STRING, got %s", args[0].Type())}
			}

			strArg := args[0].(*object.String)
			return &object.String{Value: strings.ToLower(strArg.Value)}
		},
	},

	"upper": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: fmt.Sprintf("argument must be STRING, got %s", args[0].Type())}
			}

			strArg := args[0].(*object.String)
			return &object.String{Value: strings.ToUpper(strArg.Value)}
		},
	},

	// TODO: replace, trim (family)
}
