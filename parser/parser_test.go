package parser

import (
	"fmt"
	"leonardjouve/ast"
	"leonardjouve/lexer"
	"leonardjouve/token"
	"testing"
)

func TestLetStatements(t *testing.T) {
	type LetStatementTest struct {
		input              string
		expectedIdentifier token.TokenLiteral
		expectedValue      interface{}
	}
	tests := []LetStatementTest{
		{
			input:              "let x = 5;",
			expectedIdentifier: "x",
			expectedValue:      5,
		},
		{
			input:              "let y = 10;",
			expectedIdentifier: "y",
			expectedValue:      10,
		},
		{
			input:              "let foo = x;",
			expectedIdentifier: "foo",
			expectedValue:      token.TokenLiteral("x"),
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		testParserErrors(t, parser)

		expectedStatementAmount := 1
		if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
			t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
		}

		if !testLetStatement(t, program.Statements[0], test.expectedIdentifier, test.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	type ReturnStatementTest struct {
		input    string
		expected string
	}
	tests := []ReturnStatementTest{
		{
			input:    "return 5;",
			expected: "5",
		},
		{
			input:    "return 10;",
			expected: "10",
		},
		{
			input:    "return x;",
			expected: "x",
		},
		{
			input:    "return x + y;",
			expected: "(x + y)",
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		testParserErrors(t, parser)

		expectedStatementAmount := 1
		if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
			t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
		}

		returnStatement, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("[Test] Invalid return statement: received %T, expected *ast.ReturnStatement", returnStatement)
			continue
		}

		returnType, ok := token.GetKeywordFromType(token.RETURN)
		if !ok {
			t.Error("[Test] Invalid token type: received token.RETURN")
			continue
		}

		if returnStatement.TokenLiteral() != returnType {
			t.Errorf("[Test] Invalid statement token literal: received %v, expected %v", returnStatement.TokenLiteral(), returnType)
		}

		if returnStatement.Value.String() != test.expected {
			t.Errorf("[Test] Invalid retrun expression: received %s, expected %s", returnStatement.Value.String(), test.expected)
		}
	}

}

func TestIdentifierExpressions(t *testing.T) {
	input := "foo;"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	testParserErrors(t, parser)

	expectedStatementAmount := 1
	if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
		t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
	}

	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.EpressionStatement", program.Statements[0])
	}

	if !testIdentifier(t, expressionStatement.Value, token.TokenLiteral("foo")) {
		return
	}
}

func TestIntegerExpressions(t *testing.T) {
	input := "5;"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	testParserErrors(t, parser)

	expectedStatementAmount := 1
	if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
		t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
	}

	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.EpressionStatement", program.Statements[0])
	}

	integer, ok := expressionStatement.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.IntegerLiteral", expressionStatement.Value)
	}

	var expectedIntegerValue int64 = 5
	if integer.Value != expectedIntegerValue {
		t.Errorf("[Test] Invalid identifier value: received %d, expected %d", integer.Value, expectedIntegerValue)
	}

	expectedIntegerLiteral := "5"
	if integerLiteral := string(integer.TokenLiteral()); integerLiteral != expectedIntegerLiteral {
		t.Fatalf("[Test] Invalid identifier token literal: received %s, expected %s", integerLiteral, expectedIntegerLiteral)
	}
}

func TestPrefixExpressions(t *testing.T) {
	type PrefixTest struct {
		input    string
		operator string
		value    interface{}
	}
	tests := []PrefixTest{
		{
			input:    "!5",
			operator: "!",
			value:    5,
		},
		{
			input:    "-15;",
			operator: "-",
			value:    15,
		},
		{
			input:    "!true",
			operator: "!",
			value:    true,
		},
		{
			input:    "!false",
			operator: "!",
			value:    false,
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		testParserErrors(t, parser)

		expectedStatementAmount := 1
		if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
			t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.ExpressionStatement", statement)
		}

		expression, ok := statement.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("[Test] Invalid expression type: received %T, expected *ast.PrefixExpression", expression)
		}

		if expression.Operator != test.operator {
			t.Fatalf("[Test] Invalid expression operator: received %s, expected %s", expression.Operator, test.operator)
		}

		if !testLiteralExpression(t, expression.Right, test.value) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	type InfixTest struct {
		input    string
		operator string
		left     interface{}
		right    interface{}
	}
	tests := []InfixTest{
		{
			input:    "5 + 15;",
			operator: "+",
			left:     5,
			right:    15,
		},
		{
			input:    "5 - 15;",
			operator: "-",
			left:     5,
			right:    15,
		},
		{
			input:    "5 * 15;",
			operator: "*",
			left:     5,
			right:    15,
		},
		{
			input:    "5 / 15;",
			operator: "/",
			left:     5,
			right:    15,
		},
		{
			input:    "5 < 15;",
			operator: "<",
			left:     5,
			right:    15,
		},
		{
			input:    "5 > 15;",
			operator: ">",
			left:     5,
			right:    15,
		},
		{
			input:    "5 == 15;",
			operator: "==",
			left:     5,
			right:    15,
		},
		{
			input:    "5 != 15;",
			operator: "!=",
			left:     5,
			right:    15,
		},
		{
			input:    "true == true",
			operator: "==",
			left:     true,
			right:    true,
		},
		{
			input:    "false != false",
			operator: "!=",
			left:     false,
			right:    false,
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		testParserErrors(t, parser)

		expectedStatementAmount := 1
		if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
			t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.InfixExpression", program.Statements[0])
		}

		if !testInfixExpression(t, statement.Value, test.operator, test.left, test.right) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	type PrecedenceOperatorTest struct {
		input    string
		expected string
	}
	tests := []PrecedenceOperatorTest{
		{
			input:    "-a * b",
			expected: "((-a) * b)",
		},
		{
			input:    "!-a",
			expected: "(!(-a))",
		},
		{
			input:    "a + b + c",
			expected: "((a + b) + c)",
		},
		{
			input:    "a + b - c",
			expected: "((a + b) - c)",
		},
		{
			input:    "a * b * c",
			expected: "((a * b) * c)",
		},
		{
			input:    "a * b / c",
			expected: "((a * b) / c)",
		},
		{
			input:    "a + b / c",
			expected: "(a + (b / c))",
		},
		{
			input:    "a + b * c + d / e - f",
			expected: "(((a + (b * c)) + (d / e)) - f)",
		},
		{
			input:    "3 + 4; -5 * 5",
			expected: "(3 + 4)((-5) * 5)",
		},
		{
			input:    "5 > 4 == 3 < 4",
			expected: "((5 > 4) == (3 < 4))",
		},
		{
			input:    "5 > 4 != 3 < 4",
			expected: "((5 > 4) != (3 < 4))",
		},
		{
			input:    "3 + 4 * 5 == 3 * 1 + 4 * 5",
			expected: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			input:    "true",
			expected: "true",
		},
		{
			input:    "false",
			expected: "false",
		},
		{
			input:    "3 > 5 == false",
			expected: "((3 > 5) == false)",
		},
		{
			input:    "3 < 5 == true",
			expected: "((3 < 5) == true)",
		},
		{
			input:    "1 + (2 + 3) + 4",
			expected: "((1 + (2 + 3)) + 4)",
		},
		{
			input:    "(5 + 5) * 2",
			expected: "((5 + 5) * 2)",
		},
		{
			input:    "2 / (5 + 5)",
			expected: "(2 / (5 + 5))",
		},
		{
			input:    "-(5 + 5)",
			expected: "(-(5 + 5))",
		},
		{
			input:    "!(true == true)",
			expected: "(!(true == true))",
		},
		{
			input:    "a + add(b * c) + d",
			expected: "((a + add((b * c))) + d)",
		},
		{
			input:    "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			expected: "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			input:    "add(a + b + c * d / f + g)",
			expected: "add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		testParserErrors(t, parser)

		if received := program.String(); received != test.expected {
			t.Errorf("[Test] Invalid program string: received %s, expected %s", received, test.expected)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	testParserErrors(t, parser)

	expectedStatementAmount := 1
	if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
		t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
	}

	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.EpressionStatement", program.Statements[0])
	}

	if !testBooleanLiteral(t, expressionStatement.Value, true) {
		return
	}
}

func TestIfExpressions(t *testing.T) {
	type IfExpressionTest struct {
		input             string
		conditionOperator string
		conditionLeft     token.TokenLiteral
		conditionRight    token.TokenLiteral
		consequence       token.TokenLiteral
		alternative       token.TokenLiteral
	}
	tests := []IfExpressionTest{
		{
			input:             "if (x < y) { x }",
			conditionOperator: "<",
			conditionLeft:     "x",
			conditionRight:    "y",
			consequence:       "x",
			alternative:       "",
		},
		{
			input:             "if (y > x) { y } else { x }",
			conditionOperator: ">",
			conditionLeft:     "y",
			conditionRight:    "x",
			consequence:       "y",
			alternative:       "x",
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		testParserErrors(t, parser)

		if !testIfExpression(t, program, test.conditionOperator, test.conditionLeft, test.conditionRight, test.consequence, test.alternative) {
			continue
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn (x, y) { x + y; }"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	testParserErrors(t, parser)

	expectedStatementAmount := 1
	if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
		t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
	}

	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.ExpressionStatement", program.Statements[0])
	}

	functionLiteral, ok := expressionStatement.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("[Test] Invalid expression type: received %T, expected *ast.FunctionLiteral", expressionStatement.Value)
	}

	expectedParamsAmount := 2
	if paramsAmount := len(functionLiteral.Parameters); paramsAmount != expectedParamsAmount {
		t.Fatalf("[Test] Invalid function params amount: received %d, expected %d", paramsAmount, expectedParamsAmount)
	}

	testLiteralExpression(t, functionLiteral.Parameters[0], token.TokenLiteral("x"))
	testLiteralExpression(t, functionLiteral.Parameters[1], token.TokenLiteral("y"))

	expectedBodyStatementAmount := 1
	if bodyStatementAmount := len(functionLiteral.Body.Statements); bodyStatementAmount != expectedBodyStatementAmount {
		t.Fatalf("[Test] Invalid body statement amount: received %d, expected %d", bodyStatementAmount, expectedBodyStatementAmount)
	}

	bodyStatement, ok := functionLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("[Test] Invalid body statement type: received %T, expected *ast.ExpressionStatement", functionLiteral.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Value, "+", token.TokenLiteral("x"), token.TokenLiteral("y"))
}

func TestFunctionParametersParsing(t *testing.T) {
	type FunctionParameterParsingTest struct {
		input    string
		expected []token.TokenLiteral
	}
	tests := []FunctionParameterParsingTest{
		{
			input:    "fn () {}",
			expected: []token.TokenLiteral{},
		},
		{
			input: "fn (x) {}",
			expected: []token.TokenLiteral{
				"x",
			},
		},
		{
			input: "fn (x, y) {}",
			expected: []token.TokenLiteral{
				"x",
				"y",
			},
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		testParserErrors(t, parser)

		expectedStatementAmount := 1
		if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
			t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
		}

		expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.ExpressionStatement", program.Statements[0])
		}

		functionLiteral, ok := expressionStatement.Value.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("[Test] Invalid expression type: received %T, expected *ast.FunctionLiteral", expressionStatement.Value)
		}

		expectedParamsAmount := len(test.expected)
		if paramsAmount := len(functionLiteral.Parameters); paramsAmount != expectedParamsAmount {
			t.Fatalf(" [Test] Invalid params amount: received %d, expected %d", paramsAmount, expectedParamsAmount)
		}

		for i, expectedParam := range test.expected {
			testLiteralExpression(t, functionLiteral.Parameters[i], expectedParam)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 + 3, 4 * 5);"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	testParserErrors(t, parser)

	expectedStatementAmount := 1
	if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
		t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
	}

	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.ExpressionStatement", program.Statements[0])
	}

	callExpression, ok := expressionStatement.Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf("[Test] Invalid expression type: received %T, expected *ast.CallExpression", expressionStatement.Value)
	}

	if !testIdentifier(t, callExpression.Function, token.TokenLiteral("add")) {
		return
	}

	expectedArgumentAmount := 3
	if argumentAmount := len(callExpression.Arguments); argumentAmount != expectedArgumentAmount {
		t.Fatalf("[Test] Invalid argument amount: received %d, expected %d", argumentAmount, expectedArgumentAmount)
	}

	testLiteralExpression(t, callExpression.Arguments[0], 1)
	testInfixExpression(t, callExpression.Arguments[1], "+", 2, 3)
	testInfixExpression(t, callExpression.Arguments[2], "*", 4, 5)
}

func TestCallExpressionArgumentsParsing(t *testing.T) {
	type CallExpressionArgumentsParsingTest struct {
		input              string
		expectedIdentifier token.TokenLiteral
		expectedArgs       []string
	}
	tests := []CallExpressionArgumentsParsingTest{
		{
			input:              "add()",
			expectedIdentifier: "add",
			expectedArgs:       []string{},
		},
		{
			input:              "substract(1)",
			expectedIdentifier: "substract",
			expectedArgs: []string{
				"1",
			},
		},
		{
			input:              "multiply(1, 2 + 3, 4 * 5)",
			expectedIdentifier: "multiply",
			expectedArgs: []string{
				"1",
				"(2 + 3)",
				"(4 * 5)",
			},
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		parser := New(lex)
		program := parser.ParseProgram()
		testParserErrors(t, parser)

		expectedStatementAmount := 1
		if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
			t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
		}

		expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.ExpressionStatement", program.Statements[0])
		}

		callExpression, ok := expressionStatement.Value.(*ast.CallExpression)
		if !ok {
			t.Fatalf("[Test] Invalid expression type: received %T, expected *ast.CallExpression", expressionStatement.Value)
		}

		if !testIdentifier(t, callExpression.Function, test.expectedIdentifier) {
			return
		}

		expectedArgsAmount := len(test.expectedArgs)
		if paramsAmount := len(callExpression.Arguments); paramsAmount != expectedArgsAmount {
			t.Fatalf("[Test] Invalid params amount: received %d, expected %d", paramsAmount, expectedArgsAmount)
		}

		for i, expectedArg := range test.expectedArgs {
			if arg := callExpression.Arguments[i].String(); arg != expectedArg {
				t.Errorf("[Test] Invalid call argument: received %s, expected %s", arg, expectedArg)
			}
		}
	}
}

func testParserErrors(t *testing.T, parser *Parser) {
	errorsAmount := len(parser.Errors)
	if errorsAmount == 0 {
		return
	}
	t.Errorf("[Test] Parser encountered %d error(s)", errorsAmount)
	for _, err := range parser.Errors {
		t.Error(err)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, identifier token.TokenLiteral, value interface{}) bool {
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("[Test] Invalid let statement: received %T, expected *ast.LetStatement", letStatement)
		return false
	}

	letType, ok := token.GetKeywordFromType(token.LET)
	if !ok {
		t.Error("[Test] Invalid token type: received token.LET")
		return false
	}

	if letStatement.TokenLiteral() != letType {
		t.Errorf("[Test] Invalid statement token literal: received %v, expected %v", letStatement.TokenLiteral(), letType)
		return false
	}

	if letStatement.Name.Value != identifier {
		t.Errorf("[Test] Invalid let statement name value: received %v, expected %v", letStatement.Name.Value, identifier)
		return false
	}

	if !testLiteralExpression(t, letStatement.Value, value) {
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) bool {
	integer, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("[Test] Invalid expression type: received %T, expected *ast.IntegerLiteral", integer)
		return false
	}

	if integer.Value != value {
		t.Errorf("[Test] Invalid integer value: received %d, expected %d", integer.Value, value)
		return false
	}

	expectedTokenLiteral := token.TokenLiteral(fmt.Sprintf("%d", value))
	if integer.TokenLiteral() != expectedTokenLiteral {
		t.Errorf("[Test] Invalid integer token literal: received %s, expected %s", integer.TokenLiteral(), expectedTokenLiteral)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value token.TokenLiteral) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("[Test] Invalid expression type: received %T, expected *ast.Identifier", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("[Test] Invalid identifier value: received %s, expected %s", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("[Test] Invalid identifier token literal: received %s, expected %s", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case token.TokenLiteral:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	default:
		t.Errorf("[Test] Invalid expression type: received %T", expected)
		return false
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, operator string, left interface{}, right interface{}) bool {
	infixExpression, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("[Test] Invalid expression type: received %T, expected *ast.InfixExpression", exp)
		return false
	}

	if infixExpression.Operator != operator {
		t.Errorf("[Test] invalid infix operator: received %s, expected %s", infixExpression.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, infixExpression.Left, left) {
		return false
	}

	if !testLiteralExpression(t, infixExpression.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("[Test] Invalid expression type: received %T, expected *ast.Boolean", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("[Test] Invalid identifier value: received %t, expected %t", boolean.Value, value)
		return false
	}

	expectedTokenLiteral := token.TokenLiteral(fmt.Sprintf("%t", value))
	if boolean.TokenLiteral() != expectedTokenLiteral {
		t.Errorf("[Test] Invalid identifier token literal: received %s, expected %s", boolean.TokenLiteral(), expectedTokenLiteral)
		return false
	}

	return true
}

func testIfExpression(t *testing.T, program *ast.Program, expectedConditionOperator string, expectedConditionLeft token.TokenLiteral, expecetedConditionRight token.TokenLiteral, expectedConsequence token.TokenLiteral, expectedAlternative token.TokenLiteral) bool {
	expectedStatementAmount := 1
	if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
		t.Errorf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
		return false
	}

	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("[Test] Invalid statement type: received %T, expected *ast.Expression", program.Statements[0])
		return false
	}

	ifExpression, ok := expressionStatement.Value.(*ast.IfExpression)
	if !ok {
		t.Errorf("[Test] Invalid expression type: received %T, expected *ast.IfExpression", expressionStatement)
		return false
	}

	if !testInfixExpression(t, ifExpression.Condition, expectedConditionOperator, expectedConditionLeft, expecetedConditionRight) {
		return false
	}

	expectedConsequenceStatementAmount := 1
	if consequenceStatementAmount := len(ifExpression.Consequence.Statements); consequenceStatementAmount != expectedConsequenceStatementAmount {
		t.Errorf("[Test] Invalid consequence statement amount: received %d, expected %d", consequenceStatementAmount, expectedConsequenceStatementAmount)
		return false
	}

	consequence, ok := ifExpression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("[Test] Invalid consequence statement type: received %T, expected *ast.Expression", ifExpression.Consequence.Statements[0])
		return false
	}

	if !testIdentifier(t, consequence.Value, expectedConsequence) {
		return false
	}

	if len(expectedAlternative) > 0 {
		expectedAlternativeStatementAmount := 1
		if alternativeStatementAmount := len(ifExpression.Alternative.Statements); alternativeStatementAmount != expectedAlternativeStatementAmount {
			t.Errorf("[Test] Invalid alternative statement amount: received %d, expected %d", alternativeStatementAmount, expectedAlternativeStatementAmount)
			return false
		}

		alternative, ok := ifExpression.Alternative.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("[Test] Invalid alternative statement type: received %T, expected *ast.ExpressionStatement", ifExpression.Alternative.Statements[0])
			return false
		}

		if !testIdentifier(t, alternative.Value, expectedAlternative) {
			return false
		}
	} else {
		if ifExpression.Alternative != nil {
			t.Errorf("[Test] Invalid alternative value: received %s, expected nil", ifExpression.Alternative.String())
			return false
		}
	}
	return true
}
