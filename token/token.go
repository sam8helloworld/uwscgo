package token

import "strings"

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	EOL     = "EOL"

	PROCEDURE = "PROCEDURE"
	FUNCTION  = "FUNCTION"
	RESULT    = "RESULT"
	FEND      = "FEND"

	// PRINT = "PRINT"

	// FOR  = "FOR"
	// TO   = "TO"
	// IN   = "IN"
	// NEXT = "NEXT"

	DIM    = "DIM"
	PUBLIC = "PUBLIC"
	CONST  = "CONST"

	HASHTBL = "HASHTBL"

	TRUE  = "TRUE"
	FALSE = "FALSE"

	EQUAL_OR_ASSIGN       = "="
	NOT_EQUAL             = "<>"
	LESS_THAN             = "<"
	LESS_THAN_OR_EQUAL    = "<="
	GREATER_THAN          = ">"
	GREATER_THAN_OR_EQUAL = ">="
	PLUS                  = "+"
	MINUS                 = "-"
	ASTERISK              = "*"
	SLASH                 = "/"
	MOD                   = "MOD"
	BANG                  = "!"
	LEFT_PARENTHESIS      = "("
	RIGHT_PARENTHESIS     = ")"
	LEFT_SQUARE_BRACKET   = "["
	RIGHT_SQUARE_BRACKET  = "]"
	LEFT_BRACKET          = "{"
	RIGHT_BRACKET         = "}"
	COMMA                 = ","

	IF     = "IF"
	ELSEIF = "ELSEIF"
	ELSE   = "ELSE"
	IFB    = "IFB"
	ENDIF  = "ENDIF"
	THEN   = "THEN"

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
	STRING = "STRING"
)

var reservedWords = map[string]TokenType{
	"DIM":       DIM,
	"PUBLIC":    PUBLIC,
	"CONST":     CONST,
	"HASHTBL":   HASHTBL,
	"TRUE":      TRUE,
	"FALSE":     FALSE,
	"MOD":       MOD,
	"IF":        IF,
	"ELSEIF":    ELSEIF,
	"ELSE":      ELSE,
	"IFB":       IFB,
	"ENDIF":     ENDIF,
	"THEN":      THEN,
	"FUNCTION":  FUNCTION,
	"FEND":      FEND,
	"PROCEDURE": PROCEDURE,
	"RESULT":    RESULT,
}

func LookupIdent(ident string) TokenType {
	if tokType, ok := reservedWords[strings.ToUpper(ident)]; ok {
		return tokType
	}
	return IDENT
}
