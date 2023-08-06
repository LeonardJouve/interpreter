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
	LOWERGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedence = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,
	token.LR:        LOWERGREATER,
	token.GR:        LOWERGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.ASTERISX:  PRODUCT,
	token.SLASH:     PRODUCT,
	token.LPAREN:    CALL,
}

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
		token.TRUE:       parser.parseBoolean,
		token.FALSE:      parser.parseBoolean,
		token.LPAREN:     parser.parseGroupedExpression,
		token.IF:         parser.parseIfExpression,
		token.FUNCTION:   parser.parseFunctionLiteral,
	}
}

func (parser *Parser) addInfixParsers() {
	parser.infixParsers = map[token.TokenType]infixParser{
		token.EQUAL:     parser.parseInfixExpression,
		token.NOT_EQUAL: parser.parseInfixExpression,
		token.LR:        parser.parseInfixExpression,
		token.GR:        parser.parseInfixExpression,
		token.PLUS:      parser.parseInfixExpression,
		token.MINUS:     parser.parseInfixExpression,
		token.ASTERISX:  parser.parseInfixExpression,
		token.SLASH:     parser.parseInfixExpression,
		token.LPAREN:    parser.parseCallExpression,
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

func (parser *Parser) getPrecedence() int {
	prec, ok := precedence[parser.tok.Type]
	if !ok {
		return LOWEST
	}
	return prec
}

func (parser *Parser) getNextPrecedence() int {
	prec, ok := precedence[parser.nextTok.Type]
	if !ok {
		return LOWEST
	}
	return prec
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

	parser.nextToken()
	letStatement.Value = parser.parseExpression(LOWEST)

	if parser.nextTok.Type == token.SEMICOLON {
		parser.nextToken()
	}

	return letStatement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{
		Token: parser.tok,
	}

	parser.nextToken()
	returnStatement.Value = parser.parseExpression(LOWEST)

	if parser.nextTok.Type == token.SEMICOLON {
		parser.nextToken()
	}

	return returnStatement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))
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

func (parser *Parser) parseExpression(prec int) ast.Expression {
	defer untrace(trace("parseExpression"))
	prefix, ok := parser.prefixParsers[parser.tok.Type]
	if !ok {
		parser.addInvalidPrefixError(parser.tok.Type)
		return nil
	}

	left := prefix()

	for parser.nextTok.Type != token.SEMICOLON && prec < parser.getNextPrecedence() {
		infix, ok := parser.infixParsers[parser.nextTok.Type]
		if !ok {
			return left
		}

		parser.nextToken()

		left = infix(left)
	}

	return left
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: parser.tok,
		Value: parser.tok.Literal,
	}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
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
	defer untrace(trace("parsePrefixExpression"))
	prefixExpression := &ast.PrefixExpression{
		Token:    parser.tok,
		Operator: string(parser.tok.Literal),
	}

	parser.nextToken()

	prefixExpression.Right = parser.parseExpression(PREFIX)

	return prefixExpression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	infixExpression := &ast.InfixExpression{
		Token:    parser.tok,
		Operator: string(parser.tok.Literal),
		Left:     left,
	}

	prec := parser.getPrecedence()
	parser.nextToken()
	infixExpression.Right = parser.parseExpression(prec)

	return infixExpression
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{
		Token:    parser.tok,
		Function: function,
	}

	callExpression.Arguments = parser.parseCallArguments()
	if callExpression.Arguments == nil {
		return nil
	}

	return callExpression
}

func (parser *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: parser.tok,
		Value: parser.tok.Type == token.TRUE,
	}
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()

	expression := parser.parseExpression(LOWEST)

	if !parser.expectNextTokenType(token.RPAREN) {
		return nil
	}

	return expression
}

func (parser *Parser) parseIfExpression() ast.Expression {
	ifExpression := &ast.IfExpression{
		Token: parser.tok,
	}

	if !parser.expectNextTokenType(token.LPAREN) {
		return nil
	}

	parser.nextToken()

	ifExpression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectNextTokenType(token.RPAREN) {
		return nil
	}

	if !parser.expectNextTokenType(token.LBRACE) {
		return nil
	}

	ifExpression.Consequence = parser.parseBlockStatement()

	if parser.nextTok.Type == token.ELSE {
		parser.nextToken()
		if !parser.expectNextTokenType(token.LBRACE) {
			return nil
		}
		ifExpression.Alternative = parser.parseBlockStatement()
	}

	return ifExpression
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{
		Token:      parser.tok,
		Statements: []ast.Statement{},
	}

	parser.nextToken()

	for parser.tok.Type != token.RBRACE && parser.tok.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			blockStatement.Statements = append(blockStatement.Statements, statement)
		}
		parser.nextToken()
	}

	return blockStatement
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{
		Token: parser.tok,
	}

	if !parser.expectNextTokenType(token.LPAREN) {
		return nil
	}

	functionLiteral.Parameters = parser.parseFunctionParameters()
	if functionLiteral.Parameters == nil {
		return nil
	}

	if !parser.expectNextTokenType(token.LBRACE) {
		return nil
	}

	functionLiteral.Body = parser.parseBlockStatement()

	return functionLiteral
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	parser.nextToken()

	if parser.tok.Type == token.RPAREN {
		return identifiers
	}

	identifier := &ast.Identifier{
		Token: parser.tok,
		Value: parser.tok.Literal,
	}
	identifiers = append(identifiers, identifier)

	for parser.nextTok.Type == token.COMMA {
		parser.nextToken()
		parser.nextToken()
		identifier := &ast.Identifier{
			Token: parser.tok,
			Value: parser.tok.Literal,
		}
		identifiers = append(identifiers, identifier)
	}

	if !parser.expectNextTokenType(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (parser *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if parser.nextTok.Type == token.RPAREN {
		parser.nextToken()
		return args
	}

	parser.nextToken()
	args = append(args, parser.parseExpression(LOWEST))

	for parser.nextTok.Type == token.COMMA {
		parser.nextToken()
		parser.nextToken()
		args = append(args, parser.parseExpression(LOWEST))
	}

	if !parser.expectNextTokenType(token.RPAREN) {
		return nil
	}

	return args
}
