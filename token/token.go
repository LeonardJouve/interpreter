package token

type TokenType string
type TokenLiteral string

type Token struct {
	Type    TokenType
	Literal TokenLiteral
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	ASTERISX  = "*"
	SLASH     = "/"
	EQUAL     = "=="
	NOT_EQUAL = "!="

	LR = "<"
	GR = ">"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

var keywords = map[TokenLiteral]TokenType{
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
}

func New(tokenType TokenType, tokenLiteral TokenLiteral) Token {
	return Token{
		Type:    tokenType,
		Literal: tokenLiteral,
	}
}

func LookupIdentifier(identifier TokenLiteral) TokenType {
	keyword, ok := keywords[identifier]
	if ok {
		return keyword
	}
	return IDENTIFIER
}
