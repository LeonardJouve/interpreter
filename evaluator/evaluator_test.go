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
			input:    "if (10 > 1) {if (10 > 1) {return true + false;} return 10;}",
			expected: "unknown operation: BOOLEAN + BOOLEAN",
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		err, ok := eval.(*object.Error)
		if !ok {
			t.Errorf("[Test] Invalid evaluation type: received %T, expected *object.Error", eval)
			continue
		}

		if err.Value != test.expected {
			t.Errorf("[Test] Invalid error value: received %s, expected %s", err.Value, test.expected)
			continue
		}
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	par := parser.New(lex)
	program := par.ParseProgram()

	return Eval(program)
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
