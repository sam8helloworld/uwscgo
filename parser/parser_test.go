package parser_test

import (
	"fmt"
	"testing"

	"github.com/sam8helloworld/uwscgo/ast"
	"github.com/sam8helloworld/uwscgo/lexer"
	"github.com/sam8helloworld/uwscgo/parser"
)

func TestDimStatements(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{
			"整数の変数定義",
			"DIM val = 5",
			"val",
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
			}

			stmt := program.Statements[0]
			if !testDimStatement(t, stmt, tt.expectedIdentifier) {
				return
			}

			val := stmt.(*ast.DimStatement).Value
			if !testLiteralExpression(t, val, tt.expectedValue) {
				return
			}
		})
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

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := `val`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}
	if ident.Value != "val" {
		t.Errorf("ident.Value not %s. got=%s", "val", ident.Value)
	}
	if ident.TokenLiteral() != "val" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", "val", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() not %s. got=%s", "val", literal.TokenLiteral())
	}
}

func testIntegerLiteral(
	t *testing.T,
	il ast.Expression,
	value int64,
) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIdentifier(
	t *testing.T,
	exp ast.Expression,
	value string,
) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{
			"整数同士の足し算",
			"5 + 5",
			5,
			"+",
			5,
		},
		{
			"整数同士の引き算",
			"5 - 5",
			5,
			"-",
			5,
		},
		{
			"整数同士の掛け算",
			"5 * 5",
			5,
			"*",
			5,
		},
		{
			"整数同士の割り算",
			"5 / 5",
			5,
			"/",
			5,
		},
		{
			"整数同士の余り",
			"5 MOD 5",
			5,
			"MOD",
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
			}

			if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
				return
			}

			if exp.Operator != tt.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
			}

			if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
				return
			}
		})
	}
}
