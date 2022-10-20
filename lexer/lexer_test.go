package lexer_test

import (
	"testing"

	"github.com/sam8helloworld/uwscgo/lexer"
	"github.com/sam8helloworld/uwscgo/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},`

	tests := []struct {
		expected token.Token
	}{
		{
			token.Token{
				Type:    token.ASSIGN,
				Literal: "=",
			},
		},
		{
			token.Token{
				Type:    token.PLUS,
				Literal: "+",
			},
		},
		{
			token.Token{
				Type:    token.LEFT_PARENTHESIS,
				Literal: "(",
			},
		},
		{
			token.Token{
				Type:    token.RIGHT_PARENTHESIS,
				Literal: ")",
			},
		},
		{
			token.Token{
				Type:    token.LEFT_BRACKET,
				Literal: "{",
			},
		},
		{
			token.Token{
				Type:    token.RIGHT_BRACKET,
				Literal: "}",
			},
		},
		{
			token.Token{
				Type:    token.COMMA,
				Literal: ",",
			},
		},
	}

	sut := lexer.NewLexer(input)
	for i, tt := range tests {
		tok := sut.NextToken()

		if tok.Type != tt.expected.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expected.Type, tok.Type)
		}
		if tok.Literal != tt.expected.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expected.Literal, tok.Literal)
		}
	}
}
