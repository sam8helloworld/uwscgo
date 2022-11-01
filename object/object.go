package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sam8helloworld/uwscgo/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	NULL_OBJ         = "NULL"
	BOOLEAN_OBJ      = "BOOLEAN"
	FUNCTION_OBJ     = "FUNCTION"
	ERROR_OBJ        = "ERROR"
	RESULT_VALUE_OBJ = "RESULT_VALUE"
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

	out.WriteString("FUNCTION")
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
