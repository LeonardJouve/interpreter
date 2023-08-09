package evaluator

import (
	"leonardjouve/lexer"
	"leonardjouve/object"
	"leonardjouve/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	type EvalIntegerExpressionTest struct {
		input    string
		expected int64
	}
	tests := []EvalIntegerExpressionTest{
		{
			input:    "5;",
			expected: 5,
		},
		{
			input:    "10;",
			expected: 10,
		},
		{
			input:    "-5;",
			expected: -5,
		},
		{
			input:    "-10;",
			expected: -10,
		},
		{
			input:    "5 + 5 + 5 + 5 + 10;",
			expected: 30,
		},
		{
			input:    "2 * 2 * 2 * 2 * 2;",
			expected: 32,
		},
		{
			input:    "-50 + 100 + -50;",
			expected: 0,
		},
		{
			input:    "5 * 2 + 10;",
			expected: 20,
		},
		{
			input:    "5 + 2 * 10;",
			expected: 25,
		},
		{
			input:    "20 + 2 * -10;",
			expected: 0,
		},
		{
			input:    "50 / 2 * 2 + 10;",
			expected: 60,
		},
		{
			input:    "2 * (5 + 10);",
			expected: 30,
		},
		{
			input:    "3 * 3 * 3 + 10;",
			expected: 37,
		},
		{
			input:    "3 * (3 * 3) + 10;",
			expected: 37,
		},
		{
			input:    "(5 + 10 * 2 + 15 / 3) * 2 + -10;",
			expected: 50,
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		testIntegerObject(t, eval, test.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	type EvalBooleanExpressionTest struct {
		input    string
		expected bool
	}
	tests := []EvalBooleanExpressionTest{
		{
			input:    "true;",
			expected: true,
		},
		{
			input:    "false;",
			expected: false,
		},
		{
			input:    "1 < 2;",
			expected: true,
		},
		{
			input:    "1 > 2;",
			expected: false,
		},
		{
			input:    "1 < 1;",
			expected: false,
		},
		{
			input:    "1 > 1;",
			expected: false,
		},
		{
			input:    "1 == 1;",
			expected: true,
		},
		{
			input:    "1 != 1;",
			expected: false,
		},
		{
			input:    "1 == 2;",
			expected: false,
		},
		{
			input:    "1 != 2;",
			expected: true,
		},
		{
			input:    "true == true",
			expected: true,
		},
		{
			input:    "false == false",
			expected: true,
		},
		{
			input:    "true == false",
			expected: false,
		},
		{
			input:    "false == true",
			expected: false,
		},
		{
			input:    "true != true",
			expected: false,
		},
		{
			input:    "false != false",
			expected: false,
		},
		{
			input:    "true != false",
			expected: true,
		},
		{
			input:    "false != true",
			expected: true,
		},
		{
			input:    "1 < 2 == true",
			expected: true,
		},
		{
			input:    "1 < 2 == false",
			expected: false,
		},
		{
			input:    "1 > 2 == false",
			expected: true,
		},
		{
			input:    "1 > 2 == true",
			expected: false,
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		testBooleanObject(t, eval, test.expected)
	}
}

func TestEvalBangOperator(t *testing.T) {
	type EvalBangOperatorTest struct {
		input    string
		expected bool
	}
	tests := []EvalBangOperatorTest{
		{
			input:    "!true;",
			expected: false,
		},
		{
			input:    "!false;",
			expected: true,
		},
		{
			input:    "!5;",
			expected: false,
		},
		{
			input:    "!!true;",
			expected: true,
		},
		{
			input:    "!!false;",
			expected: false,
		},
		{
			input:    "!!5;",
			expected: true,
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		testBooleanObject(t, eval, test.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	type IfElseExpressionsTest struct {
		input    string
		expected interface{}
	}
	tests := []IfElseExpressionsTest{
		{
			input:    "if (true) {10}",
			expected: 10,
		},
		{
			input:    "if (false) {10}",
			expected: nil,
		},
		{
			input:    "if (1) {10}",
			expected: 10,
		},
		{
			input:    "if (1 < 2) {10}",
			expected: 10,
		},
		{
			input:    "if (1 > 2) {10}",
			expected: nil,
		},
		{
			input:    "if (1 < 2) {10} else {5}",
			expected: 10,
		},
		{
			input:    "if (1 > 2) {10} else {5}",
			expected: 5,
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		integer, ok := test.expected.(int)
		t.Logf("%t", ok)
		if !ok {
			testNullObject(t, eval)
			continue
		}
		testIntegerObject(t, eval, int64(integer))
	}
}

func TestEvalReturnStatements(t *testing.T) {
	type EvalReturnStatementsTest struct {
		input    string
		expected int64
	}
	tests := []EvalReturnStatementsTest{
		{
			input:    "return 10;",
			expected: 10,
		},
		{
			input:    "return 10; 9;",
			expected: 10,
		},
		{
			input:    "return 2 * 5; 9;",
			expected: 10,
		},
		{
			input:    "9; return 2 * 5; 9;",
			expected: 10,
		},
		{
			input:    "9; return 2 * 5; return 9;",
			expected: 10,
		},
		{
			input:    "if (true) {if (true) {return 1;} return 2}",
			expected: 1,
		},
		{
			input:    "if (10 > 1) {if (10 > 1) {return 10;}return 1;}",
			expected: 10,
		},
		{
			input:    "let f = fn(x) {return x; x + 10;};f(10);",
			expected: 10,
		},
		{
			input:    "let f = fn(x) {let result = x + 10; return result; return 10;}; f(10);",
			expected: 20,
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		testIntegerObject(t, eval, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	type ErrorHandlingTest struct {
		input    string
		expected string
	}
	tests := []ErrorHandlingTest{
		{
			input:    "5 + true;",
			expected: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:    "5 + true; 5;",
			expected: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:    "-true;",
			expected: "unknown operation: -BOOLEAN",
		},
		{
			input:    "true + false;",
			expected: "unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			input:    "5; true + false; 5;",
			expected: "unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			input:    "if (10 > 1) {return true + false;}",
			expected: "unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			input:    "\"a\" - \"b\"",
			expected: "unknown operation: STRING - STRING",
		},
		{
			input:    "if (10 > 1) {if (10 > 1) {return true + false;} return 10;}",
			expected: "unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			input:    "foo",
			expected: "identifier not found: foo",
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		testError(t, eval, test.expected)
	}
}

func TestEvalLetStatements(t *testing.T) {
	type EvalLetStatementsTests struct {
		input    string
		expected int64
	}
	tests := []EvalLetStatementsTests{
		{
			input:    "let x = 5; x;",
			expected: 5,
		},
		{
			input:    "let x = 5 * 5; x;",
			expected: 25,
		},
		{
			input:    "let x = 5; let y = x; y;",
			expected: 5,
		},
		{
			input:    "let x = 5; let y = x; let z = x + y + 5; z;",
			expected: 15,
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		testIntegerObject(t, eval, test.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) {return x + 2;}"
	eval := testEval(input)

	function, ok := eval.(*object.Function)
	if !ok {
		t.Fatalf("[Test] Invalid evaluation type: received %T, expected *object.Function", eval)
	}

	expectedParamsAmount := 1
	if paramsAmount := len(function.Parameters); paramsAmount != expectedParamsAmount {
		t.Fatalf("[Test] Invalid parameters amout: received %d, expected %d", expectedParamsAmount, paramsAmount)
	}

	expectedFunctionParameter := "x"
	if functionParameter := function.Parameters[0].String(); functionParameter != expectedFunctionParameter {
		t.Fatalf("[Test] Invalid parameter: received %s, expected %s", functionParameter, expectedFunctionParameter)
	}

	expectedFunctionBody := "return (x + 2)"
	if functionBody := function.Body.String(); functionBody != expectedFunctionBody {
		t.Fatalf("[Test] Invalid function body: received %s, expected %s", functionBody, expectedFunctionBody)
	}
}

func TestFunctionApplication(t *testing.T) {
	type FunctionApplicationTest struct {
		input    string
		expected int64
	}
	tests := []FunctionApplicationTest{
		{
			input:    "let identify = fn(x) {x;}; identify(5);",
			expected: 5,
		},
		{
			input:    "let identify = fn(x) {return x;}; identify(5);",
			expected: 5,
		},
		{
			input:    "let double = fn(x) {x * 2;}; double(5);",
			expected: 10,
		},
		{
			input:    "let double = fn(x) {return x * 2;}; double(5);",
			expected: 10,
		},
		{
			input:    "let add = fn(x, y) {x + y;}; add(5, 5);",
			expected: 10,
		},
		{
			input:    "let add = fn(x, y) {x + y;}; add(5 + 5, add(5, 5));",
			expected: 20,
		},
		{
			input:    "fn(x) {x;}(5);",
			expected: 5,
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		testIntegerObject(t, eval, test.expected)
	}
}

func TestClosures(t *testing.T) {
	input := "let newAdder = fn(x) {return fn(y) {return x + y;};}; let x = 10; let y = 10; let addTwo = newAdder(2); addTwo(2);"

	eval := testEval(input)
	testIntegerObject(t, eval, 4)
}

func TestStringLiterals(t *testing.T) {
	input := "\"hello world\""

	eval := testEval(input)
	str, ok := eval.(*object.String)
	if !ok {
		t.Fatalf("[Test] Invalid object type: received %T, expected *object.String", eval)
	}

	if expectedStringValue := "hello world"; str.Value != expectedStringValue {
		t.Fatalf("[Test] Invalid string object value: received %s, expected %s", str.Value, expectedStringValue)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := "\"hello\" + \" \" + \"world\""

	eval := testEval(input)
	str, ok := eval.(*object.String)
	if !ok {
		t.Fatalf("[Test] Invalid object type: received %T, expected *object.String", eval)
	}

	if expectedStringValue := "hello world"; str.Value != expectedStringValue {
		t.Fatalf("[Test] Invalid string object value: received %s, expected %s", str.Value, expectedStringValue)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	type BuiltinFunctionTest struct {
		input    string
		expected interface{}
	}
	tests := []BuiltinFunctionTest{
		{
			input:    "len(\"\")",
			expected: 0,
		},
		{
			input:    "len(\"four\")",
			expected: 4,
		},
		{
			input:    "len(\"hello world\")",
			expected: 11,
		},
		{
			input:    "len(1)",
			expected: "unsupported argument for builtin function len: INTEGER",
		},
		{
			input:    "len(\"one\", \"two\")",
			expected: "wrong arguments amount: received 2, expected 1",
		},
		{
			input:    "len([1, 2, 3])",
			expected: 3,
		},
		{
			input:    "len([])",
			expected: 0,
		},
		{
			input:    "first([1, 2, 3])",
			expected: 1,
		},
		{
			input:    "first([])",
			expected: nil,
		},
		{
			input:    "first(1)",
			expected: "unsupported argument for builtin function first: INTEGER",
		},
		{
			input:    "last([1, 2, 3])",
			expected: 3,
		},
		{
			input:    "last([])",
			expected: nil,
		},
		{
			input:    "last(1)",
			expected: "unsupported argument for builtin function last: INTEGER",
		},
		{
			input:    "rest([1, 2, 3])",
			expected: []int{2, 3},
		},
		{
			input:    "rest([])",
			expected: nil,
		},
		{
			input:    "push([], 1)",
			expected: []int{1},
		},
		{
			input:    "push(1, 1)",
			expected: "unsupported argument for builtin function push: INTEGER",
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)

		switch expected := test.expected.(type) {
		case int:
			testIntegerObject(t, eval, int64(expected))
		case string:
			testError(t, eval, expected)
		case []int:
			array, ok := eval.(*object.Array)
			if !ok {
				t.Errorf("[Test] Invalid object type: received %T, expected *object.Array", eval)
				continue
			}

			expectedElementAmount := len(array.Value)
			if elementAmount := len(array.Value); elementAmount != expectedElementAmount {
				t.Errorf("[Test] Invalid array element amount: received %d, expected %d", elementAmount, expectedElementAmount)
				continue
			}

			for i, expectedElement := range expected {
				testIntegerObject(t, array.Value[i], int64(expectedElement))
			}
		case nil:
			testNullObject(t, eval)
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	eval := testEval(input)
	array, ok := eval.(*object.Array)
	if !ok {
		t.Fatalf("[Test] Invalid object type: received %T, expected *object.Array", eval)
	}

	expectedElementAmount := 3
	if elementAmount := len(array.Value); elementAmount != expectedElementAmount {
		t.Fatalf("[Test] Invalid element amount: received %d, expected %d", elementAmount, expectedElementAmount)
	}

	testIntegerObject(t, array.Value[0], 1)
	testIntegerObject(t, array.Value[1], 4)
	testIntegerObject(t, array.Value[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	type ArrayIndexExpressionTest struct {
		input    string
		expected interface{}
	}
	tests := []ArrayIndexExpressionTest{
		{
			input:    "[1, 2, 3][0];",
			expected: 1,
		},
		{
			input:    "[1, 2, 3][1];",
			expected: 2,
		},
		{
			input:    "[1, 2, 3][2];",
			expected: 3,
		},
		{
			input:    "let x = 0; [1][x];",
			expected: 1,
		},
		{
			input:    "[1, 2, 3][1 + 1];",
			expected: 3,
		},
		{
			input:    "let x = [1, 2, 3]; x[2];",
			expected: 3,
		},
		{
			input:    "let x = [1, 2, 3]; x[0] + x[1] + x[2];",
			expected: 6,
		},
		{
			input:    "let x = [1, 2, 3]; let y = x[0]; x[y];",
			expected: 2,
		},
		{
			input:    "[1, 2, 3][3];",
			expected: nil,
		},
		{
			input:    "[1, 2, 3][-1];",
			expected: nil,
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		expected, ok := test.expected.(int)
		if !ok {
			testNullObject(t, eval)
			continue
		}
		testIntegerObject(t, eval, int64(expected))
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	par := parser.New(lex)
	program := par.ParseProgram()
	env := object.NewEnvironement()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	integer, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("[Test] Invalid object type: received %T, expected *object.Integer", obj)
		return false
	}

	if integer.Value != expected {
		t.Errorf("[Test] Invalid integer value: received %d, expected %d", integer.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	boolean, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("[Test] Invalid object type: received %T, expected *object.Boolean", obj)
		return false
	}

	if boolean.Value != expected {
		t.Errorf("[Test] Invalid boolean value: received %t, expected %t", boolean.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("[Test] Invalid null object: received %T, expected: *evaluator.NULL", obj)
		return false
	}
	return true
}

func testError(t *testing.T, obj object.Object, expected string) bool {
	err, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("[Test] Invalid evaluation type: received %T, expected *object.Error", obj)
		return false
	}

	if err.Value != expected {
		t.Errorf("[Test] Invalid error value: received %s, expected %s", err.Value, expected)
		return false
	}
	return true
}
