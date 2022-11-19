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
			if len(args) >= 2 {
				array, ok := args[0].Value.(*object.Array)
				if !ok {
					return newError("argument 1 to `CALCARRAY` not supported, got %s", args[0].Value.Type())
				}
				cons, ok := args[1].Value.(*object.BuiltinConstant)
				if !ok {
					return newError("argument 2 to `CALCARRAY` not supported, got %s", args[1].Value.Type())
				}
				from := int64(0)
				to := int64(len(array.Elements) - 1)
				if len(args) >= 3 {
					fromArg, ok := args[2].Value.(*object.Integer)
					if ok {
						from = fromArg.Value
					} else {
						_, ok := args[2].Value.(*object.Empty)
						if ok {
							from = 0
						} else {
							return newError("argument 3 to `CALCARRAY` not supported, got %s", args[2].Value.Type())
						}
					}
					if len(args) == 4 {
						toArg, ok := args[3].Value.(*object.Integer)
						if !ok {
							return newError("argument 4 to `CALCARRAY` not supported, got %s", args[3].Value.Type())
						}
						to = toArg.Value
					}
				}

				if cons.T == CALC_ADD {
					var sum int64
					for i := from; i <= to; i++ {
						v, ok := array.Elements[i].(*object.Integer)
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
				if cons.T == CALC_MIN {
					var min int64
					for i := from; i <= to; i++ {
						v, ok := array.Elements[i].(*object.Integer)
						if !ok {
							return newError("array of argument 1 has not integer element. array[%d]=%q", i, v)
						}
						if i == 0 {
							min = v.Value
						}
						if v.Value < min {
							min = v.Value
						}
					}
					return &object.BuiltinFuncReturnResult{
						Value: &object.Integer{
							Value: min,
						},
					}
				}
				if cons.T == CALC_MAX {
					var max int64
					for i := from; i <= to; i++ {
						v, ok := array.Elements[i].(*object.Integer)
						if !ok {
							return newError("array of argument 1 has not integer element. array[%d]=%q", i, v)
						}
						if i == 0 {
							max = v.Value
						}
						if v.Value > max {
							max = v.Value
						}
					}
					return &object.BuiltinFuncReturnResult{
						Value: &object.Integer{
							Value: max,
						},
					}
				}
			}
			// TODO: CALC_AVRによる平均値を求める処理
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		},
	},
}

const (
	CALC_ADD       = object.BuiltinConstantType("CALC_ADD")
	CALC_MIN       = object.BuiltinConstantType("CALC_MIN")
	CALC_MAX       = object.BuiltinConstantType("CALC_MAX")
	CALC_AVR       = object.BuiltinConstantType("CALC_AVR")
	HASH_CASECARE  = object.BuiltinConstantType("HASH_CASECARE")
	HASH_SORT      = object.BuiltinConstantType("HASH_SORT")
	HASH_EXISTS    = object.BuiltinConstantType("HASH_EXISTS")
	HASH_REMOVE    = object.BuiltinConstantType("HASH_REMOVE")
	HASH_KEY       = object.BuiltinConstantType("HASH_KEY")
	HASH_VAL       = object.BuiltinConstantType("HASH_VAL")
	HASH_REMOVEALL = object.BuiltinConstantType("HASH_REMOVEALL")
)

var builtinConstants = map[object.BuiltinConstantType]object.Object{
	CALC_ADD: &object.BuiltinConstant{
		T: CALC_ADD,
		Value: &object.Integer{
			Value: 1,
		},
	},
	CALC_MIN: &object.BuiltinConstant{
		T: CALC_MIN,
		Value: &object.Integer{
			Value: 2,
		},
	},
	CALC_MAX: &object.BuiltinConstant{
		T: CALC_MAX,
		Value: &object.Integer{
			Value: 3,
		},
	},
	CALC_AVR: &object.BuiltinConstant{
		T: CALC_AVR,
		Value: &object.Integer{
			Value: 4,
		},
	},
	HASH_CASECARE: &object.BuiltinConstant{
		T: HASH_CASECARE,
		Value: &object.Integer{
			Value: 4096,
		},
	},
	HASH_SORT: &object.BuiltinConstant{
		T: HASH_SORT,
		Value: &object.Integer{
			Value: 8192,
		},
	},
	HASH_EXISTS: &object.BuiltinConstant{
		T: HASH_EXISTS,
		Value: &object.Integer{
			Value: -103,
		},
	},
	HASH_REMOVE: &object.BuiltinConstant{
		T: HASH_REMOVE,
		Value: &object.Integer{
			Value: -104,
		},
	},
	HASH_KEY: &object.BuiltinConstant{
		T: HASH_KEY,
		Value: &object.Integer{
			Value: -101,
		},
	},
	HASH_VAL: &object.BuiltinConstant{
		T: HASH_VAL,
		Value: &object.Integer{
			Value: -102,
		},
	},
	HASH_REMOVEALL: &object.BuiltinConstant{
		T: HASH_REMOVEALL,
		Value: &object.Integer{
			Value: -109,
		},
	},
}

func builtin(key string) (object.Object, bool) {
	k := strings.ToUpper(key)
	if result, ok := builtinConstants[object.BuiltinConstantType(k)]; ok {
		return result, ok
	}
	if result, ok := builtinFunctions[k]; ok {
		return result, ok
	}
	return nil, false
}
