package baisl

type TokenType int8

const (
	TokenType_UNKNOWN TokenType = iota
	TokenType_EOF
	TokenType_IDENTIFIER
	TokenType_NUMBER
	TokenType_LPAREN
	TokenType_RPAREN
	TokenType_LBRACE
	TokenType_RBRACE
	TokenType_COLON
	TokenType_COMMA
	TokenType_KEYW_FN
	TokenType_KEYW_INT
	TokenType_KEYW_VOID
	TokenType_KEYW_RETURN
)

var TokenTypeToKeyword = map[TokenType]string{
	TokenType_KEYW_FN:     "fn",
	TokenType_KEYW_VOID:   "void",
	TokenType_KEYW_INT:    "int",
	TokenType_KEYW_RETURN: "return",
}

var KeywordToTokenType = map[string]TokenType{
	"fn":     TokenType_KEYW_FN,
	"int":    TokenType_KEYW_INT,
	"void":   TokenType_KEYW_VOID,
	"return": TokenType_KEYW_RETURN,
}

func IsKeywordTokenType(tokenType TokenType) bool {
	return tokenType >= TokenType_KEYW_FN
}

func (ttype TokenType) String() string {
	switch ttype {
	case TokenType_UNKNOWN:
		return "UNKNOWN"
	case TokenType_EOF:
		return "EOF"
	case TokenType_IDENTIFIER:
		return "IDENTIFIER"
	case TokenType_NUMBER:
		return "NUMBER"
	case TokenType_LPAREN:
		return "LPAREN"
	case TokenType_RPAREN:
		return "RPAREN"
	case TokenType_LBRACE:
		return "LBRACE"
	case TokenType_RBRACE:
		return "RBRACE"
	case TokenType_COLON:
		return "COLON"
	case TokenType_KEYW_FN:
		return "KEYW_FN"
	case TokenType_KEYW_VOID:
		return "KEYW_VOID"
	case TokenType_KEYW_RETURN:
		return "KEYW_RETURN"
	case TokenType_KEYW_INT:
		return "KEYW_INT"
	default:
		return "UNKNOWN"
	}
}
