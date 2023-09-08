package evaluator

import (
	"leonardjouve/lexer"
	"leonardjouve/object"
	"leonardjouve/parser"
	"leonardjouve/token"
	"testing"
)

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1;
	let func = fn(x, y) {x + y;};
	let mac = macro(x, y) {x + y;};
	`

	lex := lexer.New(input)
	par := parser.New(lex)
	program := par.ParseProgram()
	env := object.NewEnvironement()

	DefineMacros(program, env)

	expectedStatementAmount := 2
	if statementAmount := len(program.Statements); statementAmount != expectedStatementAmount {
		t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", statementAmount, expectedStatementAmount)
	}

	_, ok := env.Get("number")
	if ok {
		t.Fatalf("[Test] Invalid environement: \"number\" is undefined")
	}

	_, ok = env.Get(token.TokenLiteral("func"))
	if ok {
		t.Fatalf("[Test] Invalid environement: \"func\" is undefined")
	}

	obj, ok := env.Get(token.TokenLiteral("mac"))
	if !ok {
		t.Fatalf("[Test] Invalid environement; \"mac\" is undefined")
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("[Test] Invalid object type: received %T, expected *object.Macro", obj)
	}

	expectedParameterAmount := 2
	if parameterAmount := len(macro.Parameters); parameterAmount != expectedParameterAmount {
		t.Fatalf("[Test] Invalid parameter amount: received %d, expected %d", parameterAmount, expectedParameterAmount)
	}

	expectedParameterString := "x"
	if parameterString := macro.Parameters[0].String(); parameterString != expectedParameterString {
		t.Fatalf("[Test] Invalid parameter string: received %s, expected %s", parameterString, expectedParameterString)
	}

	expectedParameterString = "y"
	if parameterString := macro.Parameters[1].String(); parameterString != expectedParameterString {
		t.Fatalf("[Test] Invalid parameter string: received %s, expected %s", parameterString, expectedParameterString)
	}

	expectedBody := "(x + y)"
	if body := macro.Body.String(); body != expectedBody {
		t.Fatalf("[Test] Invalid body string: received %s, expected %s", body, expectedBody)
	}
}

func TestExpandMacro(t *testing.T) {
	type ExpandMacroTest struct {
		input    string
		expected string
	}
	tests := []ExpandMacroTest{
		{
			input:    "let infixExpression = macro() {quote(1 + 2);}; infixExpression();",
			expected: "(1 + 2)",
		},
		{
			input:    "let reverse = macro(a, b) {quote(unquote(b) - unquote(a));}; reverse(1 + 1, 2 + 2);",
			expected: "(2 + 2) - (1 + 1)",
		},
		{
			input:    "let unless = macro(condition, consequence, alternative) {quote(if (!unquote(condition)) {unquote(consequence)} else {unquote(alternative)})}; unless(10 > 5, puts(\"not greater\"), puts(\"greater\"));",
			expected: "if (!(10 > 5)) {puts(\"not greater\")} else {puts(\"greater\")}",
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		par := parser.New(lex)
		program := par.ParseProgram()

		expectedLex := lexer.New(test.expected)
		expectedPar := parser.New(expectedLex)
		expectedProgram := expectedPar.ParseProgram()

		env := object.NewEnvironement()
		DefineMacros(program, env)

		expanded := ExpandMacros(program, env)

		if expanded.String() != expectedProgram.String() {
			t.Errorf("[Test] Invalid macro expansion string: received %s, expected %s", expanded.String(), expectedProgram.String())
		}
	}
}
