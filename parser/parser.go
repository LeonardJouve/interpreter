package parser

import (
	"fmt"
	"leonardjouve/ast"
	"leonardjouve/lexer"
	"leonardjouve/token"
	"strconv"
)

type Parser struct {
	lex           *lexer.Lexer
	tok           token.Token
	nextTok       token.Token
	Errors        []string
	prefixParsers map[token.TokenType]prefixParser
	infixParsers  map[token.TokenType]infixParser
}

type (
	prefixParser func() ast.Expression
	infixParser  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lex:    lex,
		Errors: []string{},
	}
	parser.nextToken()
	parser.nextToken()
	parser.addPrefixParsers()
	parser.addInfixParsers()

	return parser
}

func (parser *Parser) nextToken() {
	parser.tok = parser.nextTok
	parser.nextTok = parser.lex.NextToken()
}

func (parser *Parser) addError(err string) {
	parser.Errors = append(parser.Errors, err)
}

func (parser *Parser) addInvalidNextTokenTypeError(received token.Token, expected token.TokenType) {
	parser.addError(fmt.Sprintf("[Error] Invalid next token type: received %s %s, expected %s", received.Type, received.Literal, expected))
}

func (parser *Parser) addInvalidPrefixError(tokenType token.TokenType) {
	parser.addError(fmt.Sprintf("[Error] Invalid prefix for %s", tokenType))
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for parser.tok.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) addPrefixParsers() {
	parser.prefixParsers = map[token.TokenType]prefixParser{
		token.IDENTIFIER: parser.parseIdentifier,
		token.INT:        parser.parseIntegerLiteral,
		token.MINUS:      parser.parsePrefixExpression,
		token.BANG:       parser.parsePrefixExpression,
	}
}

func (parser *Parser) addInfixParsers() {
	parser.infixParsers = map[token.TokenType]infixParser{}
}

func (parser *Parser) addPrefixParser(tokenType token.TokenType, fn prefixParser) {
	parser.prefixParsers[tokenType] = fn
}

func (parser *Parser) addInfixParser(tokenType token.TokenType, fn infixParser) {
	parser.infixParsers[tokenType] = fn
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.tok.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) expectNextTokenType(tokenType token.TokenType) bool {
	if parser.nextTok.Type != tokenType {
		parser.addInvalidNextTokenTypeError(parser.nextTok, tokenType)
		return false
	}
	parser.nextToken()
	return true
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{
		Token: parser.tok,
	}

	if !parser.expectNextTokenType(token.IDENTIFIER) {
		return nil
	}

	letStatement.Name = &ast.Identifier{
		Token: parser.tok,
		Value: parser.tok.Literal,
	}

	if !parser.expectNextTokenType(token.ASSIGN) {
		return nil
	}

	// TODO: parse expression
	for parser.tok.Type != token.SEMICOLON {
		parser.nextToken()
	}

	return letStatement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{
		Token: parser.tok,
	}

	// TODO: parse expression
	for parser.tok.Type != token.SEMICOLON {
		parser.nextToken()
	}

	return returnStatement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionStatement := &ast.ExpressionStatement{
		Token: parser.tok,
	}

	value := parser.parseExpression(LOWEST)
	if value == nil {
		return nil
	}
	expressionStatement.Value = value

	if parser.nextTok.Type == token.SEMICOLON {
		parser.nextToken()
	}

	return expressionStatement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := parser.prefixParsers[parser.tok.Type]
	if !ok {
		parser.addInvalidPrefixError(parser.tok.Type)
		return nil
	}
	return prefix()
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: parser.tok,
		Value: parser.tok.Literal,
	}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	integer := &ast.IntegerLiteral{
		Token: parser.tok,
	}
	value, err := strconv.ParseInt(string(parser.tok.Literal), 0, 64)
	if err != nil {
		err := fmt.Sprintf("[Error] Invalid token literal. Could not parse %s as int", parser.tok.Literal)
		parser.addError(err)
		return nil
	}
	integer.Value = value

	return integer
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	prefixExpression := &ast.PrefixExpression{
		Token:    parser.tok,
		Operator: string(parser.tok.Literal),
	}

	parser.nextToken()

	prefixExpression.Right = parser.parseExpression(PREFIX)

	return prefixExpression
}
