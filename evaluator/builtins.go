package evaluator

import (
	"strings"
	"unicode/utf8"

	"github.com/sam8helloworld/uwscgo/object"
)

var builtins = map[string]*object.Builtin{
	"LENGTH": {
		Fn: func(args ...object.BuiltinFuncArgument) object.Object {
			if len(args) != 1 {
				return &object.BuiltinFuncReturnResult{
					Value: newError("argument to `LENGTH` not supported, got %s", args[0].Value.Type()),
				}
			}
			switch arg := args[0].Value.(type) {
			case *object.String:
				return &object.BuiltinFuncReturnResult{
					Value: &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))},
				}
			case *object.Array:
				return &object.BuiltinFuncReturnResult{
					Value: &object.Integer{Value: int64(len(arg.Elements))},
				}
			default:
				return &object.BuiltinFuncReturnResult{
					Value: newError("argument to `LENGTH` not supported, got %s", args[0].Value.Type()),
				}
			}
		},
	},
	"RESIZE": {
		Fn: func(args ...object.BuiltinFuncArgument) object.Object {
			if len(args) == 1 {
				array, ok := args[0].Value.(*object.Array)
				if !ok {
					return newError("argument 1 to `RESIZE` not supported, got %s", args[0].Value.Type())
				}
				return &object.BuiltinFuncReturnResult{
					Value: &object.Integer{Value: int64(len(array.Elements)) - 1},
				}
			}
			if len(args) == 2 {
				array, ok := args[0].Value.(*object.Array)
				if !ok {
					return newError("argument 1 to `RESIZE` not supported, got %s", args[0].Value.Type())
				}
				size, ok := args[1].Value.(*object.Integer)
				if !ok {
					return newError("argument 2 to `RESIZE` not supported, got %s", args[1].Value.Type())
				}

				resizedElements := make([]object.Object, len(array.Elements), size.Value+1)
				for i := 0; i < int(size.Value)+1; i++ {
					if i > len(array.Elements)-1 {
						resizedElements = append(resizedElements, EMPTY)
					} else {
						resizedElements[i] = array.Elements[i]
					}
				}
				return &object.BuiltinFuncReturnReference{
					Expression: args[0].Expression,
					Value: &object.Array{
						Elements: resizedElements,
					},
					Result: size,
				}
			}
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		},
	},
}

func builtin(key string) (*object.Builtin, bool) {
	k := strings.ToUpper(key)
	result, ok := builtins[k]
	return result, ok
}
