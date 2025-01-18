package modules

import "github.com/ajtroup1/clear/object"

var ArraysBuiltins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return &object.Error{Message: "argument to `len` not supported"}
			}
		},
	},

	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return &object.Error{Message: "first argument must be ARRAY"}
			}

			arr := args[0].(*object.Array)
			arr.Elements = append(arr.Elements, args[1:]...)

			return arr
		},
	},

	"pop": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return &object.Error{Message: "first argument must be ARRAY"}
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length == 0 {
				return &object.Null{}
			}

			popped := arr.Elements[length-1]
			arr.Elements = arr.Elements[:length-1]

			return popped
		},
	},

	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return &object.Error{Message: "first argument must be ARRAY"}
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length == 0 {
				return &object.Null{}
			}

			return arr.Elements[0]
		},
	},

	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return &object.Error{Message: "first argument must be ARRAY"}
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length == 0 {
				return &object.Null{}
			}

			newElements := make([]object.Object, length-1)
			copy(newElements, arr.Elements[1:length])

			return &object.Array{Elements: newElements}
		},
	},

	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return &object.Error{Message: "first argument must be ARRAY"}
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length == 0 {
				return &object.Null{}
			}

			return arr.Elements[length-1]
		},
	},

	"reverse": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return &object.Error{Message: "first argument must be ARRAY"}
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length == 0 {
				return arr
			}

			newElements := make([]object.Object, length)
			for i, el := range arr.Elements {
				newElements[length-i-1] = el
			}

			return &object.Array{Elements: newElements}
		},
	},
}
