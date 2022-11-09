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
			`DIM val = 0
IF TRUE THEN val = 10
val`,
			10,
		},
		{
			"条件式がFALSEの真偽値の場合THENの後の式を処理しない",
			`DIM val = 0
IF FALSE THEN val = 10
val`,
			0,
		},
		{
			"条件式が0の場合THENの後の式を処理しない",
			`DIM val = 0
IF 0 THEN val = 10
val`,
			0,
		},
		{
			"条件式が0以外の場合THENの後の式を処理する",
			`DIM val = 0
IF 1 THEN val = 10
val`,
			10,
		},
		{
			"条件式の評価がTRUEになる場合THENの後の式を処理する",
			`DIM val = 0
IF 1 < 2 THEN val = 10 ELSE val = 20
val`,
			10,
		},
		{
			"条件式の評価がFALSEになる場合ELSEの後の式を処理する",
			`DIM val = 0
IF 1 > 2 THEN val = 10 ELSE val = 20
val`,
			20,
		},
		{
			"IFBの条件式の評価がTRUEになる場合THENの後の式を処理する",
			`IFB 1 < 2 THEN
			DIM val = 10
		ELSE
			DIM val = 20
		ENDIF
			val
					`,
			10,
		},
		{
			"IFBの条件式の評価がFALSEになり、ELSEIFの条件式の評価がTRUEになる場合ELSEIFのTHENの後の式を処理する",
			`IFB 1 > 2 THEN
			DIM val = 10
		ELSEIF 1 < 2 THEN
			DIM val = 20
		ELSE
			DIM val = 30
		ENDIF
		val
		`,
			20,
		},
		{
			"IFBの条件式の評価もELSEIFの条件式の評価もFALSEになる場合ELSEの後の式を処理する",
			`IFB 1 > 2 THEN
			DIM val = 10
		ELSEIF 1 > 2 THEN
			DIM val = 20
		ELSE
		DIM val = 30
		ENDIF
		val`,
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

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expected           int64
		expectErrorMessage string
	}{
		{
			"引数なし",
			`FUNCTION fn()
	RESULT = 5
FEND
fn()`,
			5,
			"",
		},
		{
			"引数1つ",
			`FUNCTION fn(x)
	RESULT = x
FEND
fn(5)`,
			5,
			"",
		},
		{
			"引数を利用した変数をRESULTに代入",
			`FUNCTION fn(x)
	DIM y = x + 5
	RESULT = y
FEND
fn(5)`,
			10,
			"",
		},
		{
			"RESULTがない場合はエラーになる",
			`FUNCTION fn(x)
	x
FEND
fn(5)`,
			0,
			"result value does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectErrorMessage != "" {
				evaluated := testEval(tt.input)
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Fatalf("no error object returned. got=%T (%+v)", evaluated, errObj)
				}
				if errObj.Message != tt.expectErrorMessage {
					t.Errorf("wrong error message. expected=%s, got=%s", tt.expectErrorMessage, errObj.Message)
				}
			} else {
				testIntegerObject(t, testEval(tt.input), tt.expected)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedMessage string
	}{
		{
			"型が異なるもの同士の足し算",
			"5 + TRUE",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"予期せぬ前置演算子",
			"-TRUE",
			"unknown operator: -BOOLEAN",
		},
		{
			"予期せぬ中置演算子",
			"TRUE + FALSE",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"文字列同士の予期せぬ中置演算子",
			`"Hello" - "World!"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("no error object returned. got=%T (%+v)", evaluated, errObj)
			}
			if errObj.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected=%s, got=%s", tt.expectedMessage, errObj.Message)
			}
		})
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"LENGTH_空文字の時は0を返す",
			`LENGTH("")`,
			0,
		},
		{
			"LENGTH_文字がある場合は空白も含めた文字数を返す",
			`LENGTH("hello world!")`,
			12,
		},
		{
			"LENGTH_配列の場合は配列の要素数を返す",
			`DIM array[] = 1, 2, 3
LENGTH(array)`,
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				}
				if errObj.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
				}
			}
		})
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `DIM array[2] = 1, 2 * 2, 3 + 3
array
	`

	evaluated := testEval(input)
	array, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(array.Elements))
	}

	testIntegerObject(t, array.Elements[0], 1)
	testIntegerObject(t, array.Elements[1], 4)
	testIntegerObject(t, array.Elements[2], 6)
}

func TestArrayLiterals_配列の要素数が宣言と異なる(t *testing.T) {
	input := `DIM array[1] = 1, 2 * 2, 3 + 3
array
	`
	expectedMessage := "array has wrong size: [1, (2 * 2), (3 + 3)]"
	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("no error object returned. got=%T (%+v)", evaluated, errObj)
	}
	if errObj.Message != expectedMessage {
		t.Errorf("wrong error message. expected=%s, got=%s", expectedMessage, errObj.Message)
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"配列の添字に直接数字を指定する",
			`DIM array[] = 1, 2, 3
array[0]`,
			1,
		},
		{
			"配列の添字に数字を代入した識別子を指定する",
			`DIM array[] = 1, 2, 3
DIM index = 0
array[index]`,
			1,
		},
		{
			"配列の添字に数字の計算結果を指定する",
			`DIM array[] = 1, 2, 3
array[1 + 1]`,
			3,
		},
		{
			"配列のそれぞれの値を加算する",
			`DIM array[] = 1, 2, 3
array[0] + array[1] + array[2]`,
			6,
		},
		{
			"配列の存在しない要素を指定する",
			`DIM array[] = 1, 2, 3
array[3]`,
			nil,
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
