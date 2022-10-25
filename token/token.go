package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	SEMICOLON = ";"

	// PROCEDURE = "PROCEDURE"
	// FUNCTION  = "FUNCTION"
	// RESULT    = "RESULT"
	// FEND      = "FEND"

	// PRINT = "PRINT"

	// FOR  = "FOR"
	// TO   = "TO"
	// IN   = "IN"
	// NEXT = "NEXT"

	DIM    = "DIM"
	PUBLIC = "PUBLIC"

	ASSIGN            = "="
	PLUS              = "+"
	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
	LEFT_BRACKET      = "{"
	RIGHT_BRACKET     = "}"
	COMMA             = ","

	// IF     = "IF"
	// ELSEIF = "ELSEIF"
	// ELSE   = "ELSE"
	// IFB    = "IFB"
	// ENDIF  = "ENDIF"

	// WHILE = "WHILE"
	// WEND  = "WEND"

	// REPEAT = "REPEAT"
	// UNTIL  = "UNTIL"

	// BREAK    = "BREAK"
	// CONTINUE = "CONTINUE"

	// SELECT  = "SELECT"
	// CASE    = "CASE"
	// DEFAULT = "DEFAULT"
	// SELEND  = "SELEND"

	// CALL = "CALL"

	IDENT = "IDENT"
	INT   = "INT"
	// EXPANDABLE_STRING = "EXPANDABLE_STRING"
	// STRING            = "STRING"
)

var reservedWords = map[string]TokenType{
	"DIM":    DIM,
	"PUBLIC": PUBLIC,
}

func LookupIdent(ident string) TokenType {
	if tokType, ok := reservedWords[ident]; ok {
		return tokType
	}
	return IDENT
}
