package evaluator

import (
	"fmt"
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

func Eval(node ast.Node, env *object.Environement) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Value, env)
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
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
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		return &object.Return{
			Value: value,
		}
	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
		return value
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		arguments := evalExpressions(node.Arguments, env)
		if len(arguments) == 1 && isError(arguments[0]) {
			return arguments[0]
		}

		return applyFunction(function, arguments)
	case *ast.StringLiteral:
		return &object.String{
			Value: node.Value,
		}
	default:
		return nil
	}
}

func evalProgram(statements []ast.Statement, env *object.Environement) object.Object {
	var obj object.Object

	for _, statement := range statements {
		obj = Eval(statement, env)

		switch obj := obj.(type) {
		case *object.Return:
			return obj.Value
		case *object.Error:
			return obj
		}
	}

	return obj
}

func evalBlockStatement(node *ast.BlockStatement, env *object.Environement) object.Object {
	var obj object.Object

	for _, statement := range node.Statements {
		obj = Eval(statement, env)

		if objType := obj.Type(); obj != nil && (objType == object.RETURN || objType == object.ERROR) {
			return obj
		}
	}

	return obj
}

func nativeBoolToBooleanObject(boolean bool) *object.Boolean {
	if boolean {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangeOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return &object.Error{
			Value: fmt.Sprintf("unknown operation: %s%s", operator, right.Type()),
		}
	}
}

func evalBangeOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return &object.Error{
			Value: fmt.Sprintf("unknown operation: -%s", right.Type()),
		}
	}

	return &object.Integer{
		Value: -right.(*object.Integer).Value,
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return &object.Error{
			Value: fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type()),
		}
	default:
		return &object.Error{
			Value: fmt.Sprintf("unknown operation: %s %s %s", left.Type(), operator, right.Type()),
		}
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{
			Value: leftValue + rightValue,
		}

	case "-":
		return &object.Integer{
			Value: leftValue - rightValue,
		}

	case "*":
		return &object.Integer{
			Value: leftValue * rightValue,
		}
	case "/":
		return &object.Integer{
			Value: leftValue / rightValue,
		}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return &object.Error{
			Value: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()),
		}
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return &object.Error{
			Value: fmt.Sprintf("unknown operation: %s %s %s", left.Type(), operator, right.Type()),
		}
	}

	leftValue, ok := left.(*object.String)
	if !ok {
		return &object.Error{
			Value: fmt.Sprintf("invalid object type: received %T, expected *object.String", left),
		}
	}
	rightValue, ok := right.(*object.String)
	if !ok {
		return &object.Error{
			Value: fmt.Sprintf("invalid object type: received %T, expected *object.String", right),
		}
	}
	return &object.String{
		Value: leftValue.Value + rightValue.Value,
	}
}

func evalIfExpression(ifExpression *ast.IfExpression, env *object.Environement) object.Object {
	condition := Eval(ifExpression.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ifExpression.Consequence, env)
	} else if ifExpression.Alternative != nil {
		return Eval(ifExpression.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(identifier *ast.Identifier, env *object.Environement) object.Object {
	value, ok := env.Get(identifier.Value)
	if ok {
		return value
	}

	builtin, ok := builtins[identifier.Value]
	if ok {
		return builtin
	}

	return &object.Error{
		Value: fmt.Sprintf("identifier not found: %s", identifier.String()),
	}
}

func evalExpressions(expressions []ast.Expression, env *object.Environement) []object.Object {
	var exps []object.Object

	for _, exp := range expressions {
		eval := Eval(exp, env)
		if isError(eval) {
			return []object.Object{
				eval,
			}
		}
		exps = append(exps, eval)
	}

	return exps
}

func applyFunction(function object.Object, arguments []object.Object) object.Object {
	switch function := function.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnvironement(function, arguments)
		eval := Eval(function.Body, extendedEnv)
		return unwrapReturnValue(eval)
	case *object.Builtin:
		return function.Value(arguments...)
	default:
		return &object.Error{
			Value: fmt.Sprintf("not a function: %s", function.Inspect()),
		}
	}
}

func extendFunctionEnvironement(function *object.Function, arguments []object.Object) *object.Environement {
	enclosedEnv := object.NewEnclosedEnvironement(function.Env)

	for i, parameter := range function.Parameters {
		enclosedEnv.Set(parameter.Value, arguments[i])
	}

	return enclosedEnv
}

func isTruthy(obj object.Object) bool {
	return obj != FALSE && obj != NULL
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}

func unwrapReturnValue(obj object.Object) object.Object {
	returnObject, ok := obj.(*object.Return)
	if !ok {
		return obj
	}

	return returnObject.Value
}
