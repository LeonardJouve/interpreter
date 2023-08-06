package evaluator

import (
	"leonardjouve/ast"
	"leonardjouve/object"
)

var (
	FALSE = &object.Boolean{
		Value: false,
	}
	TRUE = &object.Boolean{
		Value: true,
	}
	NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return Eval(node.Statements[len(node.Statements)-1])
	case *ast.ExpressionStatement:
		return Eval(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	default:
		return nil
	}
}
