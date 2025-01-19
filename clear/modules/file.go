package modules

import (
	"os"

	"github.com/ajtroup1/clear/object"
)

var FileBuiltins = map[string]*object.Builtin{
	"read": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "argument must be STRING"}
			}

			fileName := args[0].(*object.String).Value

			data, err := os.ReadFile(fileName)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}

			return &object.String{Value: string(data)}
		},
	},

	"create": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "argument must be STRING"}
			}

			fileName := args[0].(*object.String).Value

			file, err := os.Create(fileName)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}
			defer file.Close()

			return &object.Null{}
		},
	},

	"write": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "arguments must be STRING"}
			}

			fileName := args[0].(*object.String).Value
			data := args[1].(*object.String).Value

			file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}
			defer file.Close() // Ensure the file is closed

			_, err = file.WriteString(data)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}

			return &object.Null{}
		},
	},

	"remove": &object.Builtin{
		// TODO is it a dir or file?
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "argument must be STRING"}
			}

			fileName := args[0].(*object.String).Value

			err := os.Remove(fileName)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}

			return &object.Null{}
		},
	},

	"rename": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "arguments must be STRINGs"}
			}

			oldName := args[0].(*object.String).Value
			newName := args[1].(*object.String).Value

			err := os.Rename(oldName, newName)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}

			return &object.Null{}
		},
	},

	"exists": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "argument must be STRING"}
			}

			fileName := args[0].(*object.String).Value

			_, err := os.Stat(fileName)
			if err != nil {
				return &object.Boolean{Value: false}
			}

			return &object.Boolean{Value: true}
		},
	},

	"isdir": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "argument must be STRING"}
			}

			fileName := args[0].(*object.String).Value

			info, err := os.Stat(fileName)
			if err != nil {
				return &object.Boolean{Value: false}
			}

			return &object.Boolean{Value: info.IsDir()}
		},
	},

	"isfile": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "argument must be STRING"}
			}

			fileName := args[0].(*object.String).Value

			info, err := os.Stat(fileName)
			if err != nil {
				return &object.Boolean{Value: false}
			}

			return &object.Boolean{Value: !info.IsDir()}
		},
	},
}
