package ast

import (
	"bytes"
	"leonardjouve/token"
	"strings"
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

type Boolean struct {
	Token token.Token
	Value bool
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
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

	letKeyword, ok := token.GetKeywordFromType(token.LET)
	if ok {
		out.WriteString(string(letKeyword) + " ")
	}

	out.WriteString(statement.Name.String() + " = ")

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

	returnKeyword, ok := token.GetKeywordFromType(token.RETURN)
	if ok {
		out.WriteString(string(returnKeyword) + " ")
	}

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

func (boolean *Boolean) expressionNode() {}
func (boolean *Boolean) TokenLiteral() token.TokenLiteral {
	return boolean.Token.Literal
}
func (boolean *Boolean) String() string {
	return string(boolean.Token.Literal)
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

func (statement *BlockStatement) statementNode() {}
func (statement *BlockStatement) TokenLiteral() token.TokenLiteral {
	return statement.Token.Literal
}
func (statement *BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range statement.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

func (expression *IfExpression) expressionNode() {}
func (expression *IfExpression) TokenLiteral() token.TokenLiteral {
	return expression.Token.Literal
}
func (expression *IfExpression) String() string {
	var out bytes.Buffer

	ifKeyword, ok := token.GetKeywordFromType(token.IF)
	if ok {
		out.WriteString(string(ifKeyword) + " ")
	}

	out.WriteString(expression.Condition.String() + " " + expression.Consequence.String())

	if expression.Alternative != nil {
		elseKeyword, ok := token.GetKeywordFromType(token.ELSE)
		if ok {
			out.WriteString(" " + string(elseKeyword) + " ")
		}

		out.WriteString(expression.Alternative.String())
	}

	return out.String()
}

func (expression *FunctionLiteral) expressionNode() {}
func (expression *FunctionLiteral) TokenLiteral() token.TokenLiteral {
	return expression.Token.Literal
}
func (expression *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range expression.Parameters {
		params = append(params, param.String())
	}

	functionKeyword, ok := token.GetKeywordFromType(token.FUNCTION)
	if ok {
		out.WriteString(string(functionKeyword) + " ")
	}

	out.WriteString("(" + strings.Join(params, ", ") + ") " + expression.Body.String())

	return out.String()
}

func (expression *CallExpression) expressionNode() {}
func (expression *CallExpression) TokenLiteral() token.TokenLiteral {
	return expression.Token.Literal
}
func (expression *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, argument := range expression.Arguments {
		args = append(args, argument.String())
	}

	out.WriteString(expression.Function.String() + " (" + strings.Join(args, ", ") + ")")

	return out.String()
}
