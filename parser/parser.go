package parser

import (
	"fmt"
	"leonardjouve/ast"
	"leonardjouve/lexer"
	"leonardjouve/token"
)

type Parser struct {
	lex     *lexer.Lexer
	tok     token.Token
	nextTok token.Token
	Errors  []string
}

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lex:    lex,
		Errors: []string{},
	}
	parser.nextToken()
	parser.nextToken()

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

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.tok.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return nil
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
