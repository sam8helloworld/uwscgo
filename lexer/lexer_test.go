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
			expected: token.Token{
				Type:    token.ASSIGN,
				Literal: "=",
			},
		},
		{
			expected: token.Token{
				Type:    token.PLUS,
				Literal: "+",
			},
		},
		{
			expected: token.Token{
				Type:    token.LEFT_PARENTHESIS,
				Literal: "(",
			},
		},
		{
			expected: token.Token{
				Type:    token.RIGHT_PARENTHESIS,
				Literal: ")",
			},
		},
		{
			expected: token.Token{
				Type:    token.LEFT_BRACKET,
				Literal: "{",
			},
		},
		{
			expected: token.Token{
				Type:    token.RIGHT_BRACKET,
				Literal: "}",
			},
		},
		{
			expected: token.Token{
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

func TestNextToken_整数型の変数定義(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "DIMを使った変数定義",
			input: `DIM val = 10;`,
			expected: []token.Token{
				{
					Type:    token.DIM,
					Literal: "DIM",
				},
				{
					Type:    token.IDENT,
					Literal: "val",
				},
				{
					Type:    token.ASSIGN,
					Literal: "=",
				},
				{
					Type:    token.INT,
					Literal: "10",
				},
			},
		},
		{
			name:  "PUBLICを使った変数定義",
			input: `PUBLIC val = 40;`,
			expected: []token.Token{
				{
					Type:    token.PUBLIC,
					Literal: "PUBLIC",
				},
				{
					Type:    token.IDENT,
					Literal: "val",
				},
				{
					Type:    token.ASSIGN,
					Literal: "=",
				},
				{
					Type:    token.INT,
					Literal: "40",
				},
			},
		},
	}

	for i, tt := range tests {
		sut := lexer.NewLexer(tt.input)
		for _, expected := range tt.expected {
			got := sut.NextToken()

			if got.Type != expected.Type {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, expected.Type, got.Type)
			}
			if got.Literal != expected.Literal {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, expected.Literal, got.Literal)
			}
		}
	}
}
