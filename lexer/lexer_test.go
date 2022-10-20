package lexer_test

import (
	"testing"

	"github.com/sam8helloworld/uwscgo/token"
)

func TestNextToken(t *testing.T) {
	input := `DIM val = "example_string"`

	tests := []struct {
		expected token.Token
	}{
		{
			token.Token{
				Type:    token.DIM,
				Literal: "DIM",
			},
		},
		{
			token.Token{
				Type:    token.IDENT,
				Literal: "=",
			},
		},
		{
			token.Token{
				Type:    token.ASSIGN,
				Literal: "=",
			},
		},
		{
			token.Token{
				Type:    token.EXPANDABLE_STRING,
				Literal: "example_string",
			},
		},
	}

	sut := NewLexer(input)
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
