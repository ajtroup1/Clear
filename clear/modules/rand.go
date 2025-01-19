package modules

import (
	"time"

	"math/rand"

	"github.com/ajtroup1/clear/object"
)

var RandBuiltins = map[string]*object.Builtin{
	"rand": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "wrong number of arguments"}
			}

			if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
				return &object.Error{Message: "arguments must be INTEGER"}
			}

			min := args[0].(*object.Integer).Value
			max := args[1].(*object.Integer).Value

			if min > max {
				return &object.Error{Message: "min must be less than max"}
			}

			source := rand.NewSource(time.Now().UnixNano())
			r := rand.New(source)
			randNum := r.Int63n(max-min+1) + min

			return &object.Integer{Value: randNum}
		},
	},
}
