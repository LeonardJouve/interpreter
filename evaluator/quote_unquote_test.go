package evaluator

import (
	"leonardjouve/object"
	"testing"
)

func TestQuote(t *testing.T) {
	type QuoteTest struct {
		input    string
		expected string
	}
	tests := []QuoteTest{
		{
			input:    "quote(5);",
			expected: "5",
		},
		{
			input:    "quote(5 + 8);",
			expected: "(5 + 8)",
		},
		{
			input:    "quote(foo);",
			expected: "foo",
		},
		{
			input:    "quote(foo + bar);",
			expected: "(foo + bar)",
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)

		quote, ok := eval.(*object.Quote)
		if !ok {
			anError, isAnError := eval.(*object.Error)
			if isAnError {
				t.Errorf(anError.Inspect())
			}
			t.Fatalf("[Test] Invalid object type: received %T, expected *object.Quote", eval)
		}

		if quote.Value == nil {
			t.Fatal("[Test] Invalid value: received nil")
		}

		if quote.Value.String() != test.expected {
			t.Errorf("[Test] Invalid value: received %s, expected %s", quote.Value.String(), test.expected)
		}
	}
}

func TestUnquote(t *testing.T) {
	type UnquoteTest struct {
		input    string
		expected string
	}
	tests := []UnquoteTest{
		{
			input:    "quote(unquote(4));",
			expected: "4",
		},
		{
			input:    "quote(unquote(4 + 4));",
			expected: "8",
		},
		{
			input:    "quote(unquote(4 + 4) + 4);",
			expected: "(8 + 4)",
		},
		{
			input:    "let a = 1; quote(a);",
			expected: "a",
		},
		{
			input:    "let a = 1; quote(unquote(a));",
			expected: "1",
		},
		{
			input:    "quote(unquote(true));",
			expected: "true",
		},
		{
			input:    "quote(unquote(true == false));",
			expected: "false",
		},
		{
			input:    "quote(unquote(quote(4 + 4)));",
			expected: "(4 + 4)",
		},
		{
			input:    "let quoted = quote(4 + 4); quote(unquote(4 + 4) + unquote(quoted));",
			expected: "(8 + (4 + 4))",
		},
	}

	for _, test := range tests {
		eval := testEval(test.input)

		quote, ok := eval.(*object.Quote)
		if !ok {
			anError, isAnError := eval.(*object.Error)
			if isAnError {
				t.Errorf(anError.Inspect())
			}
			t.Fatalf("[Test] Invalid object type: received %T, expected *object.Quote", eval)
		}

		if quote.Value == nil {
			t.Fatal("[Test] Invalid value: received nil")
		}

		if quote.Value.String() != test.expected {
			t.Errorf("[Test] Invalid value: received %s, expected %s", quote.Value.String(), test.expected)
		}
	}
}
