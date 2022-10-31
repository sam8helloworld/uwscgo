package evaluator_test

import (
	"testing"

	"github.com/sam8helloworld/uwscgo/evaluator"
	"github.com/sam8helloworld/uwscgo/lexer"
	"github.com/sam8helloworld/uwscgo/object"
	"github.com/sam8helloworld/uwscgo/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			"整数を評価できる",
			"5",
			5,
		},
		{
			"整数同士の足し算を評価できる",
			"5 + 5",
			10,
		},
		{
			"整数同士の引き算を評価できる",
			"5 - 5",
			0,
		},
		{
			"整数同士の掛け算を評価できる",
			"5 * 5",
			25,
		},
		{
			"整数同士の割り算を評価できる",
			"5 / 5",
			1,
		},
		{
			"整数同士の余りを評価できる",
			"5 MOD 5",
			0,
		},
		{
			"マイナスの整数を評価できる",
			"-5",
			-5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer.got=%T (%+T)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestDimStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			"整数の変数定義(1つ)",
			`DIM val = 5
val`,
			5,
		},
		{
			"整数の変数を整数の変数に代入",
			`DIM valA = 5
DIM valB = valA
valB`,
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		})
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"TRUE",
			"TRUE",
			true,
		},
		{
			"FALSE",
			"FALSE",
			false,
		},
		{
			"未満比較で左辺が右辺未満の場合はtrueを返す",
			"1 < 2",
			true,
		},
		{
			"未満比較で左辺と右辺が等しい場合はfalseを返す",
			"2 < 2",
			false,
		},
		{
			"未満比較で左辺が右辺より大きい場合はfalseを返す",
			"2 < 1",
			false,
		},
		{
			"以下比較で左辺が右辺未満の場合はtrueを返す",
			"1 <= 2",
			true,
		},
		{
			"以下比較で左辺と右辺が等しい場合はtrueを返す",
			"2 <= 2",
			true,
		},
		{
			"以下比較で左辺が右辺より大きい場合はfalseを返す",
			"2 <= 1",
			false,
		},
		{
			"超過比較で左辺が右辺未満の場合はfalseを返す",
			"1 > 2",
			false,
		},
		{
			"超過比較で左辺と右辺が等しい場合はfalseを返す",
			"2 > 2",
			false,
		},
		{
			"超過比較で左辺が右辺より大きい場合はtrueを返す",
			"2 > 1",
			true,
		},
		{
			"以上比較で左辺が右辺未満の場合はfalseを返す",
			"1 >= 2",
			false,
		},
		{
			"以上比較で左辺と右辺が等しい場合はtrueを返す",
			"2 >= 2",
			true,
		},
		{
			"以上比較で左辺が右辺より大きい場合はtrueを返す",
			"2 >= 1",
			true,
		},
		{
			"等価比較で左辺が右辺未満の場合はfalseを返す",
			"1 = 2",
			false,
		},
		{
			"等価比較で左辺と右辺が等しい場合はtrueを返す",
			"2 = 2",
			true,
		},
		{
			"等価比較で左辺が右辺より大きい場合はfalseを返す",
			"2 = 1",
			false,
		},
		{
			"等価比較の否定で左辺が右辺未満の場合はtrueを返す",
			"1 <> 2",
			true,
		},
		{
			"等価比較の否定で左辺と右辺が等しい場合はfalseを返す",
			"2 <> 2",
			false,
		},
		{
			"等価比較の否定で左辺が右辺より大きい場合はtrueを返す",
			"2 <> 1",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testBoolenObject(t, evaluated, tt.expected)
		})
	}
}

func testBoolenObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("object is not Boolean. got=%T (%+T)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestIfElseStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"条件式がTRUEの真偽値の場合THENの後の式を処理する",
			"IF TRUE THEN 10",
			10,
		},
		{
			"条件式がFALSEの真偽値の場合THENの後の式を処理しない",
			"IF FALSE THEN 10",
			nil,
		},
		{
			"条件式が0の場合THENの後の式を処理しない",
			"IF 0 THEN 10",
			nil,
		},
		{
			"条件式が0以外の場合THENの後の式を処理する",
			"IF 1 THEN 10",
			10,
		},
		{
			"条件式の評価がTRUEになる場合THENの後の式を処理する",
			"IF 1 < 2 THEN 10 ELSE 20",
			10,
		},
		{
			"条件式の評価がFALSEになる場合ELSEの後の式を処理する",
			"IF 1 > 2 THEN 10 ELSE 20",
			20,
		},
		{
			"IFBの条件式の評価がTRUEになる場合THENの後の式を処理する",
			"IFB 1 < 2 THEN 10 ELSE 20 ENDIF",
			10,
		},
		{
			"IFBの条件式の評価がFALSEになり、ELSEIFの条件式の評価がTRUEになる場合ELSEIFのTHENの後の式を処理する",
			"IFB 1 > 2 THEN 10 ELSEIF 1 < 2 THEN 20 ELSE 30 ENDIF",
			20,
		},
		{
			"IFBの条件式の評価もELSEIFの条件式の評価もFALSEになる場合ELSEの後の式を処理する",
			"IFB 1 > 2 THEN 10 ELSEIF 1 > 2 THEN 20 ELSE 30 ENDIF",
			30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
