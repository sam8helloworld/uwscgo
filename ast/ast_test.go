package ast_test

import (
	"testing"

	"github.com/sam8helloworld/uwscgo/ast"
	"github.com/sam8helloworld/uwscgo/token"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.DimStatement{
				Token: token.Token{
					Type:    token.DIM,
					Literal: "DIM",
				},
				Name: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "val",
					},
					Value: "val",
				},
				Value: &ast.Identifier{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: "5",
				},
			},
		},
	}

	if program.String() != "DIM val = 5" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
