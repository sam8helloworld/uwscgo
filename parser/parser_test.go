package parser_test

import (
	"fmt"
	"strings"
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
	case bool:
		return testBooleanLiteral(t, exp, v)
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
		leftValue  interface{}
		operator   string
		rightValue interface{}
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
		{
			"整数同士の比較(大なり)",
			"5 > 5",
			5,
			">",
			5,
		},
		{
			"整数同士の比較(小なり)",
			"5 < 5",
			5,
			"<",
			5,
		},
		{
			"整数同士の比較(等価)",
			"5 = 5",
			5,
			"=",
			5,
		},
		{
			"整数同士の比較(等価の否定)",
			"5 <> 5",
			5,
			"<>",
			5,
		},
		{
			"TRUE同士を等価比較する",
			"TRUE = TRUE",
			true,
			"=",
			true,
		},
		{
			"TRUEとFALSEを等価比較する",
			"TRUE <> FALSE",
			true,
			"<>",
			false,
		},
		{
			"FALSE同士を等価比較する",
			"FALSE = FALSE",
			false,
			"=",
			false,
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

			if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
				return
			}
		})
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestParsingPrefixExpression(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		operator     string
		integerValue int64
	}{
		{
			"エクスクラメーションマークが前置する",
			"!5",
			"!",
			5,
		},
		{
			"ハイフンが前置する(=マイナス)",
			"-15",
			"-",
			15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain %d statememts. got=%d\n", 1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}
			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
			}
			if exp.Operator != tt.operator {
				t.Fatalf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
			}

			if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
				return
			}
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"中置演算子の位置パターン01",
			"-a * b",
			"((-a) * b)",
		},
		{
			"中置演算子の位置パターン02",
			"!-a",
			"(!(-a))",
		},
		{
			"中置演算子の位置パターン03",
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"中置演算子の位置パターン04",
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"中置演算子の位置パターン05",
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"中置演算子の位置パターン06",
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"中置演算子の位置パターン07",
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"中置演算子の位置パターン08",
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"中置演算子の位置パターン09",
			`3 + 4
-5 * 5`,
			"(3 + 4)((-5) * 5)",
		},
		{
			"真偽値TRUE",
			"TRUE",
			"TRUE",
		},
		{
			"真偽値FALSE",
			"FALSE",
			"FALSE",
		},
		{
			"括弧でグルーピングするパターン01",
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"括弧でグルーピングするパターン02",
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"括弧でグルーピングするパターン03",
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"括弧でグルーピングするパターン04",
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"括弧でグルーピングするパターン05",
			"!(TRUE = TRUE)",
			"(!(TRUE = TRUE))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			actual := program.String()
			if actual != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, actual)
			}
		})
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != strings.ToUpper(fmt.Sprintf("%t", value)) {
		t.Errorf("bo.TokenLiteral not %s. got=%s", strings.ToUpper(fmt.Sprintf("%t", value)), bo.TokenLiteral())
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	input := `IF x < y THEN x ELSE y`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	// stmt := program.Statements[0]
	// if !testDimStatement(t, stmt, tt.expectedIdentifier) {
	// 	return
	// }

	// val := stmt.(*ast.DimStatement).Value
	// if !testLiteralExpression(t, val, tt.expectedValue) {
	// 	return
	// }

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	consequence, ok := exp.Consequence.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence is not ast.ExpressionStatement. got=%T", exp.Consequence)
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative is not ast.ExpressionStatement. got=%T", exp.Alternative)
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}
