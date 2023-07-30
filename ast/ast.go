package ast

import "leonardjouve/token"

type Node interface {
	TokenLiteral() token.TokenLiteral
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type Identifier struct {
	Token token.Token
	Value token.TokenLiteral
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (prog *Program) TokenLiteral() token.TokenLiteral {
	return token.TokenLiteral("")
}

func (identifier *Identifier) expressionNode() {}
func (identifier *Identifier) TokenLiteral() token.TokenLiteral {
	return identifier.Token.Literal
}

func (statement *LetStatement) statementNode() {}
func (statement *LetStatement) TokenLiteral() token.TokenLiteral {
	return statement.Token.Literal
}

func (statement *ReturnStatement) statementNode() {}
func (statement *ReturnStatement) TokenLiteral() token.TokenLiteral {
	return statement.Token.Literal
}
