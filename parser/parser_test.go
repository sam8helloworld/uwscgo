package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sam8helloworld/uwscgo/ast"
	"github.com/sam8helloworld/uwscgo/lexer"
	"github.com/sam8helloworld/uwscgo/parser"
	"github.com/sam8helloworld/uwscgo/token"
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

func TestAssignExpression(t *testing.T) {
	input := `DIM val = 0
val = 10`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	dimStmt := program.Statements[0]
	if !testDimStatement(t, dimStmt, "val") {
		return
	}

	val := dimStmt.(*ast.DimStatement).Value
	if !testLiteralExpression(t, val, 0) {
		return
	}

	exStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got=%T", exStmt)
	}

	assExp, ok := exStmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("assExp is not ast.AssignmentExpression. got=%T", assExp)
	}
	if assExp.Left.TokenLiteral() != "val" {
		t.Errorf("assStmt.Left.TokenLiteral() not '%s'. got=%s", "val", assExp.Left.TokenLiteral())
	}

	value := &ast.IntegerLiteral{
		Token: token.Token{
			Type:    token.INT,
			Literal: "10",
		},
		Value: int64(10),
	}

	rightVal, ok := assExp.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("assExp.Value is not ast.IntegerLiteral. got=%T", rightVal)
	}
	if rightVal.Value != value.Value {
		t.Errorf("rightVal.Value not %d. got=%q", value.Value, rightVal.Value)
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
		{
			"関数を使ったグルーピングのパターン01",
			"a + fn(b * c) + d",
			"((a + fn((b * c))) + d)",
		},
		{
			"関数を使ったグルーピングのパターン02",
			"fn(a, b, 1, 2 * 3, 4 + 5, fnc(6, 7 * 8))",
			"fn(a, b, 1, (2 * 3), (4 + 5), fnc(6, (7 * 8)))",
		},
		{
			"関数を使ったグルーピングのパターン03",
			"fn(a + b + c * d / f + g)",
			"fn((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"添字演算子を使ったグルーピングのパターン01",
			"a * xxx[b * c] * d",
			"((a * (xxx[(b * c)])) * d)",
		},
		{
			"添字演算子を使ったグルーピングのパターン02",
			"add(a * b[2], b[1], 2 * xxx[1])",
			"add((a * (b[2])), (b[1]), (2 * (xxx[1])))",
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

func TestIfStatement(t *testing.T) {
	input := `IF x < y THEN x ELSE y`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	consequence, ok := stmt.Consequence.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Consequence is not ast.ExpressionStatement. got=%T", stmt.Consequence)
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := stmt.Alternative.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Alternative is not ast.ExpressionStatement. got=%T", stmt.Alternative)
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestIfbStatement(t *testing.T) {
	input := `IFB x < y THEN
	y
ELSEIF x < z THEN
	z
ELSE
	x
ENDIF`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfbStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfbStatement. got=%T", program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Fatalf("stmt.Consequence.Statements is not 1 statement. got=%d\n", len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Consequence)
	}

	if !testIdentifier(t, consequence.Expression, "y") {
		return
	}

	alternative, ok := stmt.Alternative.(*ast.IfbStatement)
	if !ok {
		t.Fatalf("stmt.Alternative is not ast.IfbStatement. got=%T", stmt.Alternative)
	}

	if !testInfixExpression(t, alternative.Condition, "x", "<", "z") {
		return
	}

	if len(alternative.Consequence.Statements) != 1 {
		t.Fatalf("alternative.Consequence.Statements is not 1 statement. got=%d\n", len(alternative.Consequence.Statements))
	}

	alcons, ok := alternative.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", alternative.Consequence)
	}

	if !testIdentifier(t, alcons.Expression, "z") {
		return
	}

	el, ok := alternative.Alternative.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("alternative.Alternative is not ast.BlockStatement. got=%T", alternative.Alternative)
	}

	if len(el.Statements) != 1 {
		t.Fatalf("el.Statements is not 1 statement. got=%d\n", len(el.Statements))
	}

	st := el.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("el.Statements[0] is not ast.ExpressionStatement. got=%T", el.Statements[0])
	}
	if !testIdentifier(t, st.Expression, "x") {
		return
	}
}

func TestFunctionParsing(t *testing.T) {
	input := `FUNCTION fn(x, y)
	x + y
FEND`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	function, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T", program.Statements[0])
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function parameters wrong. want=%d, got=%d\n", 2, len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("unction.Body.Statements does not contain %d statements. got=%d\n", 1, len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedParams []string
	}{
		{
			"FUNCTION_引数なし",
			`FUNCTION fn()
	RESULT = TRUE
FEND`,
			[]string{},
		},
		{
			"FUNCTION_引数が1つ",
			`FUNCTION fn(x)
	RESULT = TRUE
FEND`,
			[]string{"x"},
		},
		{
			"FUNCTION_引数が3つ",
			`FUNCTION fn(x, y, z)
	RESULT = TRUE
FEND`,
			[]string{"x", "y", "z"},
		},
		{
			"PROCEDURE_引数なし",
			`PROCEDURE fn()
FEND`,
			[]string{},
		},
		{
			"PROCEDURE_引数が1つ",
			`PROCEDURE fn(x)
FEND`,
			[]string{"x"},
		},
		{
			"PROCEDURE_引数が3つ",
			`PROCEDURE fn(x, y, z)
FEND`,
			[]string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			funcStmt := program.Statements[0].(*ast.FunctionStatement)

			if len(funcStmt.Parameters) != len(tt.expectedParams) {
				t.Errorf("length parameters wrong. want=%d, got=%d", len(tt.expectedParams), len(funcStmt.Parameters))
			}

			for i, ident := range tt.expectedParams {
				testLiteralExpression(t, funcStmt.Parameters[i], ident)
			}
		})
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "fn(1, 2 * 3, 4 + 5)"

	l := lexer.NewLexer(input)
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

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "fn") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestResultStatements(t *testing.T) {
	input := `FUNCTION fn()
	RESULT = 5
FEND`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("stmt.Body.Statements does not contain 1 statement. got=%d", len(stmt.Body.Statements))
	}

	rstmt, ok := stmt.Body.Statements[0].(*ast.ResultStatement)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not ast.ResultStatement. got=%T", stmt.Body.Statements[0])
	}

	testIntegerLiteral(t, rstmt.ResultValue, 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}

}

func TestParsingArrayLiterals(t *testing.T) {
	input := `DIM array[2] = 1, 2 * 2, 3 + 3`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.DimStatement)
	if !ok {
		t.Fatalf("stmt not ast.DimStatement. got=%T", program.Statements[0])
	}
	array, ok := stmt.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Value)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Index, 2)
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingArrayLiterals_空配列(t *testing.T) {
	input := `DIM array[] = 1, 2 * 2, 3 + 3`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.DimStatement)
	if !ok {
		t.Fatalf("stmt not ast.DimStatement. got=%T", program.Statements[0])
	}
	array, ok := stmt.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Value)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	if array.Index != nil {
		t.Errorf("array.Index has wrong value. got=%q, want=nil", array.Index)
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingArrayLiterals_配列変数の宣言のみ(t *testing.T) {
	input := `DIM array[2]`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.DimStatement)
	if !ok {
		t.Fatalf("stmt not ast.DimStatement. got=%T", program.Statements[0])
	}
	array, ok := stmt.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Value)
	}

	if len(array.Elements) != 0 {
		t.Fatalf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Index, 2)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := `arr[1 + 1]`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "arr") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingArrayAssign(t *testing.T) {
	input := `arr[1] = 2`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	assExp, ok := stmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("exp not ast.AssignmentExpression. got=%T", stmt.Expression)
	}

	leftExp, ok := assExp.Left.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("assExp.Left not ast.IndexExpression. got=%T", assExp.Left)
	}

	index, ok := leftExp.Index.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("leftExp.Index not ast.IntegerLiteral. got=%T", leftExp.Index)
	}
	if index.Value != 1 {
		t.Errorf("index.Value is not 1. got=%d", index.Value)
	}

	value, ok := assExp.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("assExp.Value not ast.IntegerLiteral. got=%T", assExp.Value)
	}

	if value.Value != 2 {
		t.Errorf("value.Value is not 1. got=%d", index.Value)
	}
}
