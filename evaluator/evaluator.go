package evaluator

import (
	"fmt"

	"github.com/sam8helloworld/uwscgo/ast"
	"github.com/sam8helloworld/uwscgo/object"
)

var (
	NULL  = &object.Null{}
	EMPTY = &object.Empty{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.DimStatement:
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val)
	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		env.SetConst(node.Name.Value, val)
	case *ast.HashTableStatement:
		val := Eval(node.Value, env)
		evalHashTableStatement(node.Name.Value, val, env)
	case *ast.ForToStepStatement:
		evalForToStepStatement(node, env)
	case *ast.ForInStatement:
		evalForInStatement(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.EmptyArgument:
		return EMPTY
	case *ast.Empty:
		return EMPTY
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IfStatement:
		return evalIfExpression(node, env)
	case *ast.IfbStatement:
		return evalIfbStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.AssignmentExpression:
		left := node.Left
		val := Eval(node.Value, env)
		return evalAssignExpression(left, val, env)
	case *ast.ResultStatement:
		val := Eval(node.ResultValue, env)
		return &object.ResultValue{
			Value: val,
		}
	case *ast.FunctionStatement:
		name := node.Name.Value
		params := node.Parameters
		body := node.Body
		function := &object.Function{
			Name:       name,
			Parameters: params,
			Body:       body,
			Env:        env,
			IsProc:     node.IsProc,
		}
		env.Set(name, function)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		switch fn := function.(type) {
		case *object.Function:
			return applyFunction(fn, args)
		case *object.BuiltinFunction:
			argss := []object.BuiltinFuncArgument{}
			for i, arg := range args {
				argss = append(argss, object.BuiltinFuncArgument{
					Expression: node.Arguments[i],
					Value:      arg,
				})
			}
			return applyBuiltinFunction(fn, argss, env)
		default:
			return newError("not a function: %s", fn.Type())
		}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		size := Eval(node.Size, env)
		// ?????????????????????
		if len(node.Elements) == 0 {
			if size == nil {
				return &object.Array{Elements: []object.Object{}}
			}
			sizeObj, ok := Eval(node.Size, env).(*object.Integer)
			if !ok {
				return newError("array has wrong size: %s", node.String())
			}
			return &object.Array{Elements: make([]object.Object, sizeObj.Value+1)}
		}
		// ????????????????????????
		if size != nil {
			sizeObj, ok := Eval(node.Size, env).(*object.Integer)
			if !ok {
				return newError("array has wrong size: %s", node.String())
			}
			if len(node.Elements) != int(sizeObj.Value)+1 {
				return newError("array has wrong size: %s", node.String())
			}
		}
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		var opt object.Object
		if node.Option != nil {
			opt = Eval(node.Option, env)
		}
		return evalIndexExpression(left, index, opt)
	}
	return nil
}

func applyFunction(fn *object.Function, args []object.Object) object.Object {
	extendedEnv := extendFunctionEnv(fn, args)
	evaluated := Eval(fn.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func applyBuiltinFunction(fn *object.BuiltinFunction, args []object.BuiltinFuncArgument, env *object.Environment) object.Object {
	result := fn.Fn(args...)
	switch r := result.(type) {
	case *object.BuiltinFuncReturnResult:
		return r.Value
	case *object.BuiltinFuncReturnReference:
		evalAssignExpression(r.Expression, r.Value, env)
		return r.Result
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIndex, param := range fn.Parameters {
		env.Set(param.Value, args[paramIndex])
	}

	if !fn.IsProc {
		env.Set("RESULT", NULL)
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	// ?????????NULL?????????Result??????????????????????????????
	if resultValue, ok := obj.(*object.ResultValue); ok {
		return resultValue.Value
	}

	return newError("result value does not exist")
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ResultValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil && result.Type() == object.RESULT_VALUE_OBJ {
			return result
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtin(node.Value); ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{
			Value: leftVal + rightVal,
		}
	case "-":
		return &object.Integer{
			Value: leftVal - rightVal,
		}
	case "*":
		return &object.Integer{
			Value: leftVal * rightVal,
		}
	case "/":
		return &object.Integer{
			Value: leftVal / rightVal,
		}
	case "MOD":
		return &object.Integer{
			Value: leftVal % rightVal,
		}
	case "=":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "<>":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(ie *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

// NOTE: evalIfExpression???????????????????????????????????????????????????????????????
func evalIfbStatement(is *ast.IfbStatement, env *object.Environment) object.Object {
	condition := Eval(is.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(is.Consequence, env)
	} else if is.Alternative != nil {
		return Eval(is.Alternative, env)
	} else {
		return NULL
	}
}

// ??????????????????0?????????FALSE???????????????????????????????????????
func isTruthy(obj object.Object) bool {
	if obj == FALSE {
		return false
	}
	if obj.Type() == object.INTEGER_OBJ {
		v := obj.(*object.Integer).Value
		if v == 0 {
			return false
		}
	}
	return true
}

func evalAssignExpression(left ast.Expression, val object.Object, env *object.Environment) object.Object {
	switch l := left.(type) {
	case *ast.Identifier:
		ident, ok := env.Get(l.Value)
		if !ok {
			return newError("identifier is not defined: %s", l.String())
		}
		if _, ok := ident.(*object.HashTable); ok {
			if cons, ok := val.(*object.BuiltinConstant); ok {
				if cons.T == HASH_REMOVEALL {
					ht := &object.HashTable{
						Pairs: map[object.HashKey]object.HashPair{},
					}
					env.Set(l.Value, ht)
					return ht
				}
			}
		}
		env.Set(l.Value, val)
	case *ast.IndexExpression:
		ident, ok := l.Left.(*ast.Identifier)
		if !ok {
			return newError("index expression left should be identifier: %s", l.Left.String())
		}
		aoh, ok := env.Get(ident.Value)
		if !ok {
			return newError("identifier is not defined: %s", ident.String())
		}
		switch aoh := aoh.(type) {
		case *object.Array:
			index, ok := l.Index.(*ast.IntegerLiteral)
			if !ok {
				return newError("index sholud be integer: %s", l.Index.String())
			}
			aoh.Elements[int(index.Value)] = val
			env.Set(ident.Value, aoh)
		case *object.HashTable:
			index := Eval(l.Index, env)
			key, ok := index.(object.Hashable)
			if !ok {
				return newError("unusable as hash key: %s", index.Type())
			}
			aoh.Pairs[key.HashKey()] = object.HashPair{
				Key:   index,
				Value: val,
			}
			return aoh
		default:
			return newError("index should be with array or hash: %s", l.Left.String())
		}
	}
	return val
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIndexExpression(left, index, opt object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASHTBL_OBJ:
		return evalHashTableIndexExpression(left, index, opt)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

func evalHashTableIndexExpression(hash, index, opt object.Object) object.Object {
	hashObject := hash.(*object.HashTable)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	if opt != nil {
		opt, ok := opt.(*object.BuiltinConstant)
		if !ok {
			return newError("option should be builtin constant: %s", opt.Value.Inspect())
		}
		if opt.T == HASH_EXISTS {
			_, ok := hashObject.Pairs[key.HashKey()]
			return &object.Boolean{
				Value: ok,
			}
		}
		if opt.T == HASH_REMOVE {
			delete(hashObject.Pairs, key.HashKey())
		}
		if opt.T == HASH_KEY {
			i, ok := key.(*object.Integer)
			if !ok {
				return newError("unusable as hash key: %s", key.HashKey().Type)
			}
			pair := hashObject.GetPairByIndex(int(i.Value))

			return pair.Key
		}
		if opt.T == HASH_VAL {
			i, ok := key.(*object.Integer)
			if !ok {
				return newError("unusable as hash key: %s", key.HashKey().Type)
			}
			pair := hashObject.GetPairByIndex(int(i.Value))

			return pair.Value
		}
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalHashTableStatement(name string, value object.Object, env *object.Environment) object.Object {
	var casecare bool = false
	var sort bool = false
	switch val := value.(type) {
	case *object.BuiltinConstant:
		if val.T == HASH_CASECARE {
			casecare = true
		}
		if val.T == HASH_SORT {
			sort = true
		}
	case *object.Empty:
	default:
		return newError("unknown hash declare: %s", val.Inspect())
	}
	return env.Set(name, &object.HashTable{
		Pairs:      map[object.HashKey]object.HashPair{},
		IsSort:     sort,
		IsCasecare: casecare,
	})
}

func evalForToStepStatement(forStmt *ast.ForToStepStatement, env *object.Environment) object.Object {
	from, ok := forStmt.From.(*ast.IntegerLiteral)
	if !ok {
		return newError("forStmt.From is not *ast.IntegerLiteral. got=%T", forStmt.From)
	}
	to, ok := forStmt.To.(*ast.IntegerLiteral)
	if !ok {
		return newError("forStmt.To is not *ast.IntegerLiteral. got=%T", forStmt.To)
	}
	step, ok := forStmt.Step.(*ast.IntegerLiteral)
	if !ok {
		// NOTE: STEP??????????????????????????????1?????????
		step = &ast.IntegerLiteral{
			Value: int64(1),
		}
	}
Loop:
	for i := from.Value; i <= to.Value; i += step.Value {
		index := &object.Integer{
			Value: i,
		}
		env.Set(forStmt.LoopVar.Value, index)
		for _, stmt := range forStmt.Block.Statements {
			if _, ok := stmt.(*ast.ContinueStatement); ok {
				continue Loop
			}
			if _, ok := stmt.(*ast.BreakStatement); ok {
				break Loop
			}
			Eval(stmt, env)
		}
	}
	return nil
}

func evalForInStatement(forStmt *ast.ForInStatement, env *object.Environment) object.Object {
	collectionIdent, ok := forStmt.Collection.(*ast.Identifier)
	if !ok {
		return newError("forStmt.Collection is not *ast.Identifier. got=%T", forStmt.Collection)
	}

	collect := Eval(collectionIdent, env)
	collectObject, ok := collect.(*object.Array)
	if !ok {
		return newError("collect is not *object.Array. got=%T", collect)
	}
Loop:
	for _, element := range collectObject.Elements {
		env.Set(forStmt.LoopVar.Value, element)

		for _, stmt := range forStmt.Block.Statements {
			if _, ok := stmt.(*ast.ContinueStatement); ok {
				continue Loop
			}
			if _, ok := stmt.(*ast.BreakStatement); ok {
				break Loop
			}
			Eval(stmt, env)
		}
	}
	return nil
}
