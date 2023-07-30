package parser

import (
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
