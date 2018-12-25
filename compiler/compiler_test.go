package compiler

import (
	"fmt"
	"github.com/masa-suzu/monkey/ast"
	"github.com/masa-suzu/monkey/code"
	"github.com/masa-suzu/monkey/lexer"
	"github.com/masa-suzu/monkey/object"
	"github.com/masa-suzu/monkey/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1+2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Add),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1-2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Sub),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1*2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Mul),
				code.Make(code.Pop),
			},
		},
		{
			input:             "2/1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Div),
				code.Make(code.Pop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Minus),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.True),
				code.Make(code.Pop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.False),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.GreaterThan),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.GreaterThan),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Equal),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.NotEqual),
				code.Make(code.Pop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.True),
				code.Make(code.False),
				code.Make(code.Equal),
				code.Make(code.Pop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.True),
				code.Make(code.False),
				code.Make(code.NotEqual),
				code.Make(code.Pop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.True),
				code.Make(code.Bang),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"monkey"`,
			expectedConstants: []interface{}{"monkey"},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Pop),
			},
		},
		{
			input:             `"mon" + "key"`,
			expectedConstants: []interface{}{"mon", "key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Add),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "if(true){10}3333",
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{
				code.Make(code.True),
				code.Make(code.JumpNotTruthy, 10),
				code.Make(code.Constant, 0),
				code.Make(code.Jump, 11),
				code.Make(code.Null),
				code.Make(code.Pop),
				code.Make(code.Constant, 1),
				code.Make(code.Pop),
			},
		},
		{
			input:             "if(true){10}else{20}3333",
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				code.Make(code.True),
				code.Make(code.JumpNotTruthy, 10),
				code.Make(code.Constant, 0),
				code.Make(code.Jump, 13),
				code.Make(code.Constant, 1),
				code.Make(code.Pop),
				code.Make(code.Constant, 2),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let one = 1;
			let two = 2;
`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.Constant, 1),
				code.Make(code.SetGlobal, 1),
			},
		},
		{
			input: `
			let one = 1;
			one;
`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.GetGlobal, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
			let one = 1;
			let two = one;
			two;
`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.GetGlobal, 0),
				code.Make(code.SetGlobal, 1),
				code.Make(code.GetGlobal, 1),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `[]`,
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.Array, 0),
				code.Make(code.Pop),
			},
		},
		{
			input:             `[1,2,3]`,
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Constant, 2),
				code.Make(code.Array, 3),
				code.Make(code.Pop),
			},
		},
		{

			input:             `[1+2,3-4,5*6]`,
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Add),
				code.Make(code.Constant, 2),
				code.Make(code.Constant, 3),
				code.Make(code.Sub),
				code.Make(code.Constant, 4),
				code.Make(code.Constant, 5),
				code.Make(code.Mul),
				code.Make(code.Array, 3),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `{}`,
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.Hash, 0),
				code.Make(code.Pop),
			},
		},
		{
			input:             `{1:2,3:4,5:6}`,
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Constant, 2),
				code.Make(code.Constant, 3),
				code.Make(code.Constant, 4),
				code.Make(code.Constant, 5),
				code.Make(code.Hash, 6),
				code.Make(code.Pop),
			},
		},
		{

			input:             `{1:2+3}`,
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.Constant, 0),
				code.Make(code.Constant, 1),
				code.Make(code.Constant, 2),
				code.Make(code.Add),
				code.Make(code.Hash, 2),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func runCompilerTest(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compile error: %v", err)
		}
		byteCode := compiler.ByteCode()

		err = testInstructions(tt.expectedInstructions, byteCode.Instructions)

		if err != nil {
			t.Fatalf("testInstructions failed: %v", err)
		}

		err = testConstants(tt.expectedConstants, byteCode.Constants)

		if err != nil {
			t.Fatalf("testConstants failed: %v", err)
		}
	}
}

func testInstructions(
	expected []code.Instructions,
	actual code.Instructions,
) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrogn instruction at %d.\nwant=%q\ngot =%q", i, concatted, actual)
		}
	}
	return nil
}

func testConstants(
	expected []interface{},
	actual []object.Object,
) error {

	if len(expected) != len(actual) {
		return fmt.Errorf("wrong numbers of constants.\nwant=%q\ngot =%q", len(expected), len(actual))
	}

	for i, c := range expected {
		switch constant := c.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		case string:
			err := testStringObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject faild : %s", i, err)
			}
		}
	}
	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testIntegerObject(expected int64, actual object.Object) error {
	ret, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer.want=%d,got=%d", expected, actual)
	}
	if ret.Value != expected {
		return fmt.Errorf("object has wrong value. want=%d,got=%d", expected, ret.Value)
	}
	return nil
}

func testStringObject(expected string, actual object.Object) error {
	ret, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not string.want=%q,got=%d", expected, actual)
	}
	if ret.Value != expected {
		return fmt.Errorf("object has wrong value. want=%q,got=%q", expected, ret.Value)
	}
	return nil
}
