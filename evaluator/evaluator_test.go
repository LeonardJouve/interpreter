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
