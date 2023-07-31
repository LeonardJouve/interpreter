package ast

import (
	"leonardjouve/token"
	"testing"
)

func TestString(t *testing.T) {
	letType, ok := token.GetKeywordFromType(token.LET)
	if !ok {
		t.Fatal("[Test] Invalid token type: received token.LET")
	}

	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: letType,
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "variable",
					},
					Value: "variable",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "anotherVariable",
					},
					Value: "anotherVariable",
				},
			},
		},
	}
	test := "let variable = anotherVariable;"

	if program.String() != test {
		t.Errorf("[Test] Invalid program string: received %s, expected %s", program.String(), test)
	}
}
