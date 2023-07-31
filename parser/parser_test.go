package parser

import (
	"fmt"
	"leonardjouve/ast"
	"leonardjouve/lexer"
	"leonardjouve/token"
	"testing"
)

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

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foo = 12121212;
	`
	tests := []token.TokenLiteral{
		"x",
		"y",
		"foo",
	}

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	testParserErrors(t, parser)

	expectedStatementAmount := 3
	if len(program.Statements) != expectedStatementAmount {
		t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", len(program.Statements), expectedStatementAmount)
	}

	for i, test := range tests {
		if !testLetStatement(t, program.Statements[i], test) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name token.TokenLiteral) bool {
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

	if letStatement.Name.Value != name {
		t.Errorf("[Test] Invalid  let statement name value: received %v, expected %v", letStatement.Name.Value, name)
		return false
	}

	// TODO: test expression

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 12121212;
	`

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	testParserErrors(t, parser)

	expectedStatementAmount := 3
	if len(program.Statements) != expectedStatementAmount {
		t.Fatalf("[Test] Invalid statement amount: received %d, expected %d", len(program.Statements), expectedStatementAmount)
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
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

		// TODO: test expression
	}
}

func TestIdentifierExpression(t *testing.T) {
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

	identifier, ok := expressionStatement.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("[Test] Invalid statement type: received %T, expected *ast.Identifier", expressionStatement.Value)
	}

	expectedIdentifierValue := token.TokenLiteral("foo")
	if identifier.Value != expectedIdentifierValue {
		t.Errorf("[Test] Invalid identifier value: received %s, expected %s", identifier.Value, expectedIdentifierValue)
	}

	if identifier.TokenLiteral() != expectedIdentifierValue {
		t.Fatalf("[Test] Invalid identifier token literal: received %s, expected %s", identifier.TokenLiteral(), expectedIdentifierValue)
	}
}

func TestIntegerExpression(t *testing.T) {
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
		value    int64
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

		if !testIntegerLiteral(t, expression.Right, test.value) {
			return
		}
	}
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
		t.Errorf("[Error] Invalid integer token literal: received %s, expected %d", integer.TokenLiteral(), expectedTokenLiteral)
		return false
	}

	return true
}
