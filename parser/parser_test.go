package parser_test

import (
	"testing"

	"github.com/sam8helloworld/uwscgo/ast"
	"github.com/sam8helloworld/uwscgo/lexer"
	"github.com/sam8helloworld/uwscgo/parser"
)

func TestDimStatements(t *testing.T) {
	input := `
DIM valA = 5
DIM valB = 6
DIM valC = 7
	`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"valA"},
		{"valB"},
		{"valC"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testDimStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testDimStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "DIM" {
		t.Errorf("s.TokenLiteral() not 'DIM'. got=%q", s.TokenLiteral())
		return false
	}

	dimStmt, ok := s.(*ast.DimStatement)
	if !ok {
		t.Errorf("s not *ast.DimStatement. got=%T", s)
		return false
	}

	if dimStmt.Name.Value != name {
		t.Errorf("dimStmt.Name.Value not '%s'. got=%s", name, dimStmt.Name.Value)
		return false
	}

	if dimStmt.Name.TokenLiteral() != name {
		t.Errorf("dimStmt.Name.TokenLiteral() not '%s'. got=%s", name, dimStmt.Name.TokenLiteral())
		return false
	}

	return true
}
