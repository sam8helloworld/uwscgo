package evaluator

import (
	"strings"
	"unicode/utf8"

	"github.com/sam8helloworld/uwscgo/object"
)

var builtinFunctions = map[string]*object.BuiltinFunction{
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
	"CALCARRAY": {
		Fn: func(args ...object.BuiltinFuncArgument) object.Object {
			if len(args) == 2 {
				array, ok := args[0].Value.(*object.Array)
				if !ok {
					return newError("argument 1 to `CALCARRAY` not supported, got %s", args[0].Value.Type())
				}
				cons, ok := args[1].Value.(*object.BuiltinConstant)
				if !ok {
					return newError("argument 2 to `CALCARRAY` not supported, got %s", args[1].Value.Type())
				}
				if cons.Value.(*object.String).Value == "CALC_ADD" {
					var sum int64
					for i, el := range array.Elements {
						v, ok := el.(*object.Integer)
						if !ok {
							return newError("array of argument 1 has not integer element. array[%d]=%q", i, v)
						}
						sum += v.Value
					}
					return &object.BuiltinFuncReturnResult{
						Value: &object.Integer{
							Value: sum,
						},
					}
				}
			}
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		},
	},
}

var builtinConstants = map[string]object.Object{
	"CALC_ADD": &object.BuiltinConstant{
		Value: &object.String{
			Value: "CALC_ADD",
		},
	},
	"CALC_MIN": &object.BuiltinConstant{
		Value: &object.String{
			Value: "CALC_MIN",
		},
	},
	"CALC_MAX": &object.BuiltinConstant{
		Value: &object.String{
			Value: "CALC_MAX",
		},
	},
	"CALC_AVR": &object.BuiltinConstant{
		Value: &object.String{
			Value: "CALC_AVR",
		},
	},
}

func builtin(key string) (object.Object, bool) {
	k := strings.ToUpper(key)
	if result, ok := builtinConstants[k]; ok {
		return result, ok
	}
	if result, ok := builtinFunctions[k]; ok {
		return result, ok
	}
	return nil, false
}
