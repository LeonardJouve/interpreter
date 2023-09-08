package evaluator

import (
	"fmt"
	"leonardjouve/ast"
	"leonardjouve/object"
	"leonardjouve/token"
)

func quote(node ast.Node, env *object.Environement) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{
		Value: node,
	}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environement) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToAstNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	return ok && callExpression.Function.TokenLiteral() == "unquote"
}

func convertObjectToAstNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Quote:
		return obj.Value
	case *object.Integer:
		return &ast.IntegerLiteral{
			Token: token.Token{
				Type:    token.INT,
				Literal: token.TokenLiteral(fmt.Sprintf("%d", obj.Value)),
			},
			Value: obj.Value,
		}
	case *object.Boolean:
		var tok token.Token
		if obj.Value {
			tok = token.Token{
				Type:    token.TRUE,
				Literal: "true",
			}
		} else {
			tok = token.Token{
				Type:    token.FALSE,
				Literal: "false",
			}
		}
		return &ast.Boolean{
			Token: tok,
			Value: obj.Value,
		}
	default:
		return nil
	}
}
