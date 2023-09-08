package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression {
		return &IntegerLiteral{
			Value: 1,
		}
	}
	two := func() Expression {
		return &IntegerLiteral{
			Value: 2,
		}
	}

	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2

		return integer
	}

	type ModifyTest struct {
		input    Node
		expected Node
	}
	tests := []ModifyTest{
		{
			input:    one(),
			expected: two(),
		},
		{
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Value: one()},
				},
			},
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Value: two()},
				},
			},
		},
		{
			&InfixExpression{
				Left:     one(),
				Right:    two(),
				Operator: "+",
			},
			&InfixExpression{
				Left:     two(),
				Right:    two(),
				Operator: "+",
			},
		},
		{
			&InfixExpression{
				Left:     two(),
				Right:    one(),
				Operator: "+",
			},
			&InfixExpression{
				Left:     two(),
				Right:    two(),
				Operator: "+",
			},
		},
		{
			&PrefixExpression{
				Right:    one(),
				Operator: "-",
			},
			&PrefixExpression{
				Right:    two(),
				Operator: "-",
			},
		},
		{
			&IndexExpression{
				Left:  one(),
				Index: one(),
			},
			&IndexExpression{
				Left:  two(),
				Index: two(),
			},
		},
		{
			&IfExpression{
				Condition: one(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Value: one()},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Value: one()},
					},
				},
			},
			&IfExpression{
				Condition: two(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Value: two()},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Value: two()},
					},
				},
			},
		},
		{
			&ReturnStatement{
				Value: one(),
			},
			&ReturnStatement{
				Value: two(),
			},
		},
		{
			&LetStatement{
				Value: one(),
			},
			&LetStatement{
				Value: two(),
			},
		},
		{
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Value: one()},
					},
				},
			},
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Value: two()},
					},
				},
			},
		},
		{
			&ArrayLiteral{
				Value: []Expression{
					one(),
					one(),
				},
			},
			&ArrayLiteral{
				Value: []Expression{
					two(),
					two(),
				},
			},
		},
	}

	for _, test := range tests {
		modified := Modify(test.input, turnOneIntoTwo)

		if equal := reflect.DeepEqual(modified, test.expected); !equal {
			t.Errorf("[Test] Invalid deep comparaison: received %#v, expected %#v", modified, test.expected)
		}

	}

	hashLiteral := &HashLiteral{
		Value: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}

	Modify(hashLiteral, turnOneIntoTwo)

	for key, value := range hashLiteral.Value {
		expectedKey := 2
		if key, _ := key.(*IntegerLiteral); key.Value != int64(expectedKey) {
			t.Errorf("[Test] Invalid hashLiteral key: recevied %d, expected %d", key.Value, expectedKey)
		}

		expectedValue := 2
		if value, _ := value.(*IntegerLiteral); value.Value != int64(expectedValue) {
			t.Errorf("[Test] Invalid hashLiteral value: received %d, expected %d", value.Value, expectedValue)
		}
	}
}
