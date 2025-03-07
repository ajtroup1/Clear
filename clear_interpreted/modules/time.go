package modules

import (
	"time"

	"github.com/ajtroup1/clear/object"
)

var TimeBuiltins = map[string]*object.Builtin{
	"now": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return &object.String{Value: time.Now().String()}
		},
	},
}
