package ast

import (
	"bytes"
	"leonardjouve/token"
)

type Node interface {
	TokenLiteral() token.TokenLiteral
	String() string
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

type ExpressionStatement struct {
	Token token.Token
	Value Expression
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

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (program *Program) TokenLiteral() token.TokenLiteral {
	return token.TokenLiteral("")
}
func (program *Program) String() string {
	var out bytes.Buffer
	for _, statement := range program.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

func (identifier *Identifier) expressionNode() {}
func (identifier *Identifier) TokenLiteral() token.TokenLiteral {
	return identifier.Token.Literal
}
func (identifier *Identifier) String() string {
	return string(identifier.Value)
}

func (statement *ExpressionStatement) statementNode() {}
func (statement *ExpressionStatement) TokenLiteral() token.TokenLiteral {
	return statement.Token.Literal
}
func (statement *ExpressionStatement) String() string {
	if statement.Value == nil {
		return ""
	}
	return statement.Value.String()
}

func (statement *LetStatement) statementNode() {}
func (statement *LetStatement) TokenLiteral() token.TokenLiteral {
	return statement.Token.Literal
}
func (statement *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(string(statement.TokenLiteral()) + " " + statement.Name.String() + " = ")

	if statement.Value != nil {
		out.WriteString(statement.Value.String())
	}

	out.WriteByte(';')

	return out.String()
}

func (statement *ReturnStatement) statementNode() {}
func (statement *ReturnStatement) TokenLiteral() token.TokenLiteral {
	return statement.Token.Literal
}
func (statement *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(string(statement.TokenLiteral()) + " ")

	if statement.Value != nil {
		out.WriteString(statement.Value.String())
	}

	out.WriteByte(';')

	return out.String()
}

func (literal *IntegerLiteral) expressionNode() {}
func (literal *IntegerLiteral) TokenLiteral() token.TokenLiteral {
	return literal.Token.Literal
}
func (literal *IntegerLiteral) String() string {
	return string(literal.Token.Literal)
}

func (expression *PrefixExpression) expressionNode() {}
func (expression *PrefixExpression) TokenLiteral() token.TokenLiteral {
	return expression.Token.Literal
}
func (expression *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(" + expression.Operator + expression.Right.String() + ")")

	return out.String()
}

func (expression *InfixExpression) expressionNode() {}
func (expression *InfixExpression) TokenLiteral() token.TokenLiteral {
	return expression.Token.Literal
}
func (expression *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(" + expression.Left.String() + " " + expression.Operator + " " + expression.Right.String() + ")")

	return out.String()
}
