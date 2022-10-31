package lexer_test

import (
	"testing"

	"github.com/sam8helloworld/uwscgo/lexer"
	"github.com/sam8helloworld/uwscgo/token"
)

type Args struct {
	name     string
	input    string
	expected []token.Token
}

func testToken(t *testing.T, tests []Args) {
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}

func TestNextToken_四則演算(t *testing.T) {
	tests := []Args{
		{
			name:  "整数の足し算",
			input: `5 + 5`,
			expected: []token.Token{
				{
					Type:    token.INT,
					Literal: "5",
				},
				{
					Type:    token.PLUS,
					Literal: "+",
				},
				{
					Type:    token.INT,
					Literal: "5",
				},
			},
		},
		{
			name:  "整数同士の足し算",
			input: `5 + 5`,
			expected: []token.Token{
				{
					Type:    token.INT,
					Literal: "5",
				},
				{
					Type:    token.PLUS,
					Literal: "+",
				},
				{
					Type:    token.INT,
					Literal: "5",
				},
			},
		},
		{
			name:  "整数同士の引き算",
			input: `5 - 5`,
			expected: []token.Token{
				{
					Type:    token.INT,
					Literal: "5",
				},
				{
					Type:    token.MINUS,
					Literal: "-",
				},
				{
					Type:    token.INT,
					Literal: "5",
				},
			},
		},
		{
			name:  "整数同士の掛け算",
			input: `5 * 5`,
			expected: []token.Token{
				{
					Type:    token.INT,
					Literal: "5",
				},
				{
					Type:    token.ASTERISK,
					Literal: "*",
				},
				{
					Type:    token.INT,
					Literal: "5",
				},
			},
		},
		{
			name:  "整数同士の割り算",
			input: `5 / 5`,
			expected: []token.Token{
				{
					Type:    token.INT,
					Literal: "5",
				},
				{
					Type:    token.SLASH,
					Literal: "/",
				},
				{
					Type:    token.INT,
					Literal: "5",
				},
			},
		},
		{
			name:  "整数同士の余り",
			input: `5 MOD 5`,
			expected: []token.Token{
				{
					Type:    token.INT,
					Literal: "5",
				},
				{
					Type:    token.MOD,
					Literal: "MOD",
				},
				{
					Type:    token.INT,
					Literal: "5",
				},
			},
		},
	}

	testToken(t, tests)
}

func TestNextToken_整数型の変数定義(t *testing.T) {
	tests := []Args{
		{
			name:  "DIMを使った変数定義",
			input: `DIM val = 10`,
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
					Type:    token.EQUAL_OR_ASSIGN,
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
			input: `PUBLIC val = 40`,
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
					Type:    token.EQUAL_OR_ASSIGN,
					Literal: "=",
				},
				{
					Type:    token.INT,
					Literal: "40",
				},
			},
		},
	}

	testToken(t, tests)
}

func TestNextToken_真偽値(t *testing.T) {
	tests := []Args{
		{
			name:  "TRUE",
			input: `TRUE`,
			expected: []token.Token{
				{
					Type:    token.TRUE,
					Literal: "TRUE",
				},
			},
		},
		{
			name:  "FALSE",
			input: `FALSE`,
			expected: []token.Token{
				{
					Type:    token.FALSE,
					Literal: "FALSE",
				},
			},
		},
	}

	testToken(t, tests)
}

func TestNextToken_比較演算子(t *testing.T) {
	tests := []Args{
		{
			name:  "等価比較",
			input: `=`,
			expected: []token.Token{
				{
					Type:    token.EQUAL_OR_ASSIGN,
					Literal: "=",
				},
			},
		},
		{
			name:  "等価比較の否定",
			input: `<>`,
			expected: []token.Token{
				{
					Type:    token.NOT_EQUAL,
					Literal: "<>",
				},
			},
		},
		{
			name:  "未満",
			input: `<`,
			expected: []token.Token{
				{
					Type:    token.LESS_THAN,
					Literal: "<",
				},
			},
		},
		{
			name:  "以下",
			input: `<=`,
			expected: []token.Token{
				{
					Type:    token.LESS_THAN_OR_EQUAL,
					Literal: "<=",
				},
			},
		},
		{
			name:  "超過",
			input: `>`,
			expected: []token.Token{
				{
					Type:    token.GREATER_THAN,
					Literal: ">",
				},
			},
		},
		{
			name:  "以上",
			input: `>=`,
			expected: []token.Token{
				{
					Type:    token.GREATER_THAN_OR_EQUAL,
					Literal: ">=",
				},
			},
		},
	}

	testToken(t, tests)
}

func TestNextToken_条件分岐(t *testing.T) {
	tests := []Args{
		{
			name:  "IF",
			input: `IF`,
			expected: []token.Token{
				{
					Type:    token.IF,
					Literal: "IF",
				},
			},
		},
		{
			name:  "ELSEIF",
			input: `ELSEIF`,
			expected: []token.Token{
				{
					Type:    token.ELSEIF,
					Literal: "ELSEIF",
				},
			},
		},
		{
			name:  "ELSE",
			input: `ELSE`,
			expected: []token.Token{
				{
					Type:    token.ELSE,
					Literal: "ELSE",
				},
			},
		},
		{
			name:  "IFB",
			input: `IFB`,
			expected: []token.Token{
				{
					Type:    token.IFB,
					Literal: "IFB",
				},
			},
		},
		{
			name:  "ENDIF",
			input: `ENDIF`,
			expected: []token.Token{
				{
					Type:    token.ENDIF,
					Literal: "ENDIF",
				},
			},
		},
		{
			name:  "THEN",
			input: `THEN`,
			expected: []token.Token{
				{
					Type:    token.THEN,
					Literal: "THEN",
				},
			},
		},
	}

	testToken(t, tests)
}

func TestNextToken_関数(t *testing.T) {
	tests := []Args{
		{
			name:  "FUNCTION",
			input: `FUNCTION`,
			expected: []token.Token{
				{
					Type:    token.FUNCTION,
					Literal: "FUNCTION",
				},
			},
		},
		{
			name:  "FEND",
			input: `FEND`,
			expected: []token.Token{
				{
					Type:    token.FEND,
					Literal: "FEND",
				},
			},
		},
		{
			name:  "PROCEDURE",
			input: `PROCEDURE`,
			expected: []token.Token{
				{
					Type:    token.PROCEDURE,
					Literal: "PROCEDURE",
				},
			},
		},
	}

	testToken(t, tests)
}
