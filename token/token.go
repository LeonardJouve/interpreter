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
	STRING     = "STRING"

	ASSIGN    = "ASSIGN"
	PLUS      = "PLUS"
	MINUS     = "MINUS"
	BANG      = "BANG"
	ASTERISX  = "ASTERISX"
	SLASH     = "SLASH"
	EQUAL     = "EQUAL"
	NOT_EQUAL = "NOT_EQUAL"

	LR = "LR"
	GR = "GR"

	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN   = "LPAREN"
	RPAREN   = "RPAREN"
	LBRACE   = "LBRACE"
	RBRACE   = "RBRACE"
	LBRACKET = "LBRACKET"
	RBRACKET = "RRACKET"

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

func GetKeywordFromType(tokenType TokenType) (TokenLiteral, bool) {
	for keyword, tokType := range keywords {
		if tokenType == tokType {
			return keyword, true
		}
	}
	return TokenLiteral(""), false
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
