package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/sam8helloworld/uwscgo/ast"
)

type ObjectType string

const (
	INTEGER_OBJ                       = "INTEGER"
	NULL_OBJ                          = "NULL"
	EMPTY_OBJ                         = "EMPTY"
	BOOLEAN_OBJ                       = "BOOLEAN"
	FUNCTION_OBJ                      = "FUNCTION"
	ERROR_OBJ                         = "ERROR"
	RESULT_VALUE_OBJ                  = "RESULT_VALUE"
	STRING_OBJ                        = "STRING"
	BUILTIN_FUNCTION_OBJ              = "BUILTIN_FUNCTION_OBJ"
	BUILTIN_CONSTANT_OBJ              = "BUILTIN_CONSTANT_OBJ"
	ARRAY_OBJ                         = "ARRAY"
	HASHTBL_OBJ                       = "HASHTBL_OBJ"
	BUILTIN_FUNC_RETURN_RESULT_OBJ    = "BUILTIN_FUNC_RETURN_RESULT"
	BUILTIN_FUNC_RETURN_REFERENCE_OBJ = "BUILTIN_FUNC_RETURN_REFERENCE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

type Empty struct{}

func (e *Empty) Inspect() string {
	return "EMPTY"
}

func (e *Empty) Type() ObjectType {
	return EMPTY_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

type Function struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	IsProc     bool
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	if f.IsProc {
		out.WriteString("PROCEDURE")
	} else {
		out.WriteString("FUNCTION")
	}
	out.WriteString(" ")
	out.WriteString(f.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(" \n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

type ResultValue struct {
	Value Object
}

func (rv *ResultValue) Type() ObjectType {
	return RESULT_VALUE_OBJ
}

func (rv *ResultValue) Inspect() string {
	return rv.Value.Inspect()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

type BuiltinFuncArgument struct {
	Expression ast.Expression
	Value      Object
}

type BuiltinFunction struct {
	Fn func(args ...BuiltinFuncArgument) Object
}

func (bf *BuiltinFunction) Type() ObjectType {
	return BUILTIN_FUNCTION_OBJ
}

func (bf *BuiltinFunction) Inspect() string {
	return "builtin function"
}

type BuiltinConstantType string

type BuiltinConstant struct {
	T     BuiltinConstantType
	Value Object
}

func (bc *BuiltinConstant) Type() ObjectType {
	return BUILTIN_CONSTANT_OBJ
}

func (bc *BuiltinConstant) Inspect() string {
	return "builtin constant"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type BuiltinFuncReturnResult struct {
	Value Object
}

func (b *BuiltinFuncReturnResult) Type() ObjectType {
	return BUILTIN_FUNC_RETURN_RESULT_OBJ
}

func (b *BuiltinFuncReturnResult) Inspect() string {
	return b.Value.Inspect()
}

type BuiltinFuncReturnReference struct {
	Expression ast.Expression
	Value      Object
	Result     Object
}

func (b *BuiltinFuncReturnReference) Type() ObjectType {
	return BUILTIN_FUNC_RETURN_REFERENCE_OBJ
}

func (b *BuiltinFuncReturnReference) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")
	out.WriteString("(")
	out.WriteString(b.Expression.String())
	out.WriteString("=")
	out.WriteString(b.Value.Inspect())
	out.WriteString("),")
	out.WriteString("result")
	out.WriteString("=")
	out.WriteString(b.Result.Inspect())
	out.WriteString("}")
	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hashable interface {
	HashKey() HashKey
}

type HashTable struct {
	Pairs    map[HashKey]HashPair
	Sort     bool
	Casecare bool
}

func (ht *HashTable) Type() ObjectType {
	return HASHTBL_OBJ
}

func (ht *HashTable) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range ht.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
