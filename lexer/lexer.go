package lexer

import (
	"leonardjouve/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func New(input string) *Lexer {
	return &Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
		char:         0,
	}
}

func (lexer *Lexer) NextToken() token.Token {
	lexer.readChar()
	lexer.skipWhitespace()

	var tokenType token.TokenType
	tokenLiteral := token.TokenLiteral(lexer.char)

	switch lexer.char {
	case '+':
		tokenType = token.PLUS
	case '(':
		tokenType = token.LPAREN
	case ')':
		tokenType = token.RPAREN
	case '{':
		tokenType = token.LBRACE
	case '}':
		tokenType = token.RBRACE
	case ',':
		tokenType = token.COMMA
	case ';':
		tokenType = token.SEMICOLON
	case '-':
		tokenType = token.MINUS
	case '/':
		tokenType = token.SLASH
	case '*':
		tokenType = token.ASTERISX
	case '<':
		tokenType = token.LR
	case '>':
		tokenType = token.GR
	case '=':
		if nextChar := lexer.getNextChar(); nextChar == '=' {
			tokenType = token.EQUAL
			tokenLiteral += token.TokenLiteral(nextChar)
			lexer.readChar()
		} else {
			tokenType = token.ASSIGN
		}
	case '!':
		if nextChar := lexer.getNextChar(); nextChar == '=' {
			tokenType = token.NOT_EQUAL
			tokenLiteral += token.TokenLiteral(nextChar)
			lexer.readChar()
		} else {
			tokenType = token.BANG
		}
	case '"':
		tokenType = token.STRING
		tokenLiteral = lexer.readString()
	case '[':
		tokenType = token.LBRACKET
	case ']':
		tokenType = token.RBRACKET
	case 0:
		tokenType = token.EOF
	default:
		if isLetter(lexer.char) {
			tokenLiteral = lexer.readIdentifier()
			tokenType = token.LookupIdentifier(tokenLiteral)
		} else if isDigit(lexer.char) {
			tokenType = token.INT
			tokenLiteral = lexer.readNumber()
		} else {
			tokenType = token.ILLEGAL
		}
	}

	return token.New(tokenType, tokenLiteral)
}

func (lexer *Lexer) readIdentifier() token.TokenLiteral {
	position := lexer.position
	for isLetter(lexer.char) && isLetter(lexer.getNextChar()) {
		lexer.readChar()
	}
	return token.TokenLiteral(lexer.input[position:lexer.readPosition])
}

func (lexer *Lexer) readNumber() token.TokenLiteral {
	position := lexer.position
	for isDigit(lexer.char) && isDigit(lexer.getNextChar()) {
		lexer.readChar()
	}
	return token.TokenLiteral(lexer.input[position:lexer.readPosition])
}

func (lexer *Lexer) readString() token.TokenLiteral {
	position := lexer.position + 1
	for {
		lexer.readChar()
		if lexer.char == '"' || lexer.char == 0 {
			break
		}
	}

	return token.TokenLiteral(lexer.input[position:lexer.position])
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.getNextChar()
	}
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) getNextChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.char == ' ' || lexer.char == '\n' || lexer.char == '\r' || lexer.char == '\t' {
		lexer.readChar()
	}
}

func isLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
