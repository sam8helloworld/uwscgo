package evaluator

import (
	"unicode/utf8"

	"github.com/sam8helloworld/uwscgo/object"
)

var builtins = map[string]*object.Builtin{
	"LENGTH": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `LENGTH` not supported, got %s", args[0].Type())
			}
		},
	},
}
