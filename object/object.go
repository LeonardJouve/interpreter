package object

import (
	"bytes"
	"fmt"
	"leonardjouve/ast"
	"leonardjouve/token"
	"strings"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type Null struct{}

type Return struct {
	Value Object
}

type Error struct {
	Value string
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environement
}

type String struct {
	Value string
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Value BuiltinFunction
}

const (
	NULL     = "NULL"
	INTEGER  = "INTEGER"
	BOOLEAN  = "BOOLEAN"
	RETURN   = "RETURN"
	ERROR    = "ERROR"
	FUNCTION = "FUNCTION"
	STRING   = "STRING"
	BUILTIN  = "BUILTIN"
)

func (integer *Integer) Type() ObjectType {
	return INTEGER
}
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

func (boolean *Boolean) Type() ObjectType {
	return BOOLEAN
}
func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}

func (null *Null) Type() ObjectType {
	return NULL
}
func (null *Null) Inspect() string {
	return "null"
}

func (ret *Return) Type() ObjectType {
	return RETURN
}
func (ret *Return) Inspect() string {
	return ret.Value.Inspect()
}

func (err *Error) Type() ObjectType {
	return ERROR
}
func (err *Error) Inspect() string {
	return "[Error] " + err.Value
}

func (function *Function) Type() ObjectType {
	return FUNCTION
}
func (function *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range function.Parameters {
		params = append(params, param.String())
	}

	functionKeyword, ok := token.GetKeywordFromType(token.FUNCTION)
	if ok {
		out.WriteString(string(functionKeyword))
	}

	out.WriteString("(" + strings.Join(params, ", ") + ") {\n" + function.Body.String() + "\n}")

	return out.String()
}

func (str *String) Type() ObjectType {
	return STRING
}
func (str *String) Inspect() string {
	return str.Value
}

func (builtin *Builtin) Type() ObjectType {
	return BUILTIN
}
func (builtin *Builtin) Inspect() string {
	return "builtin function"
}
