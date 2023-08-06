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
			input:    "true",
			expected: true,
		},
		{
			input:    "false",
			expected: false,
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
		t.Errorf("[Test] Invalid object type: received %T, expected *obejct.Boolean", obj)
		return false
	}

	if boolean.Value != expected {
		t.Errorf("[Test] Invalid boolean value: received %t, expected %t", boolean.Value, expected)
		return false
	}
	return true
}
