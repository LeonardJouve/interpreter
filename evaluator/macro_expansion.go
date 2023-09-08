package evaluator

import (
	"fmt"
	"leonardjouve/ast"
	"leonardjouve/object"
)

func DefineMacros(program *ast.Program, env *object.Environement) {
	definitions := []int{}
	for i, statement := range program.Statements {
		if isMacroDefinition(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIndex := definitions[i]
		program.Statements = append(program.Statements[:definitionIndex], program.Statements[definitionIndex+1:]...)
	}
}

func ExpandMacros(program ast.Node, env *object.Environement) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		macro := isMacroCall(callExpression, env)
		if macro == nil {
			return node
		}

		arguments := quoteArguments(callExpression.Arguments)
		env := extendMacroEnvironement(macro, arguments)

		evaluated := Eval(macro.Body, env)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic(fmt.Sprintf("Unsupported return value in macro expansion: received %T, expected *object.Quote", evaluated))
		}

		return quote.Value
	})
}

func isMacroDefinition(statement ast.Statement) bool {
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}

func addMacro(statement ast.Statement, env *object.Environement) {
	letStatement, _ := statement.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Body:       macroLiteral.Body,
		Env:        env,
	}

	env.Set(letStatement.Name.Value, macro)
}

func isMacroCall(expression *ast.CallExpression, env *object.Environement) *object.Macro {
	identifier, ok := expression.Function.(*ast.Identifier)
	if !ok {
		return nil
	}

	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil
	}

	return macro
}

func quoteArguments(arguments []ast.Expression) []*object.Quote {
	args := []*object.Quote{}
	for _, arg := range arguments {
		args = append(args, &object.Quote{Value: arg})
	}

	return args
}

func extendMacroEnvironement(macro *object.Macro, arguments []*object.Quote) *object.Environement {
	extended := object.NewEnclosedEnvironement(macro.Env)

	for i, param := range macro.Parameters {
		extended.Set(param.Value, arguments[i])
	}

	return extended
}
