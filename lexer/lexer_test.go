package lexer

import (
	"leonardjouve/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;
	
	let add = fn(x, y) {
		x + y;
	};
	
	let result = add(five, ten);
	
	!-/*5;
	5 < 10 > 5;
	
	if (5 < 10) {
		return true;
	} else {
		return false;
	}
	
	10 == 10;
	10 != 9;
	`
	tests := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "five"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "ten"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "add"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.FUNCTION, Literal: "fn"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "x"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENTIFIER, Literal: "y"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.IDENTIFIER, Literal: "x"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENTIFIER, Literal: "y"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "result"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.IDENTIFIER, Literal: "add"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "five"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENTIFIER, Literal: "ten"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.BANG, Literal: "!"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.SLASH, Literal: "/"},
		{Type: token.ASTERISX, Literal: "*"},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.INT, Literal: "5"},
		{Type: token.LR, Literal: "<"},
		{Type: token.INT, Literal: "10"},
		{Type: token.GR, Literal: ">"},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IF, Literal: "if"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.INT, Literal: "5"},
		{Type: token.LR, Literal: "<"},
		{Type: token.INT, Literal: "10"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.INT, Literal: "10"},
		{Type: token.EQUAL, Literal: "=="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.INT, Literal: "10"},
		{Type: token.NOT_EQUAL, Literal: "!="},
		{Type: token.INT, Literal: "9"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: "\x00"},
	}
	lexer := New(input)

	for i, test := range tests {
		tok := lexer.nextToken()

		if tok.Type != test.Type {
			t.Fatalf("[Test %d] Invalid token type: received %q %q, expected %q %q", i, tok.Type, tok.Literal, test.Type, test.Literal)
		}

		if tok.Literal != test.Literal {
			t.Fatalf("[Test %d] Invalid token literal: received %q %q, expected %q %q", i, tok.Type, tok.Literal, test.Type, test.Literal)
		}
	}
}
