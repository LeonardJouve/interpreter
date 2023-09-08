package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"leonardjouve/ast"
	"leonardjouve/token"
	"strings"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
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

type Array struct {
	Value []Object
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Value map[HashKey]HashPair
}

type Quote struct {
	Value ast.Node
}

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environement
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
	ARRAY    = "ARRAY"
	HASH     = "HASH"
	QUOTE    = "QUOTE"
	MACRO    = "MACRO"
)

func (integer *Integer) Type() ObjectType {
	return INTEGER
}
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}
func (integer *Integer) HashKey() HashKey {
	return HashKey{
		Type:  integer.Type(),
		Value: uint64(integer.Value),
	}
}

func (boolean *Boolean) Type() ObjectType {
	return BOOLEAN
}
func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}
func (boolean *Boolean) HashKey() HashKey {
	var value uint64

	if boolean.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{
		Type:  boolean.Type(),
		Value: value,
	}
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
func (str *String) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(str.Value))

	return HashKey{
		Type:  str.Type(),
		Value: hash.Sum64(),
	}
}

func (builtin *Builtin) Type() ObjectType {
	return BUILTIN
}
func (builtin *Builtin) Inspect() string {
	return "builtin function"
}

func (array *Array) Type() ObjectType {
	return ARRAY
}
func (array *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range array.Value {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[" + strings.Join(elements, ", ") + "]")

	return out.String()
}

func (hash *Hash) Type() ObjectType {
	return HASH
}
func (hash *Hash) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, value := range hash.Value {
		elements = append(elements, value.Key.Inspect()+": "+value.Value.Inspect())
	}

	out.WriteString("{" + strings.Join(elements, ", ") + "}")

	return out.String()
}

func (quote *Quote) Type() ObjectType {
	return QUOTE
}
func (quote *Quote) Inspect() string {
	return "QUOTE(" + quote.Value.String() + ")"
}

func (macro *Macro) Type() ObjectType {
	return MACRO
}
func (macro *Macro) Inspect() string {
	parameters := []string{}
	for _, parameter := range macro.Parameters {
		parameters = append(parameters, parameter.String())
	}

	return "macro(" + strings.Join(parameters, ", ") + ") {\n" + macro.Body.String() + "\n}"
}
