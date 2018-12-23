package vm

import (
	"fmt"
	"github.com/masa-suzu/monkey/ast"
	"github.com/masa-suzu/monkey/compiler"
	"github.com/masa-suzu/monkey/lexer"
	"github.com/masa-suzu/monkey/object"
	"github.com/masa-suzu/monkey/parser"
	"testing"
)

type testCase struct {
	in   string
	want interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []testCase{
		{"1", 1},
		{"1 + 2", 3},
	}
	testRun(t, tests)
}

func testRun(t *testing.T, tests []testCase) {
	t.Helper()

	for _, tt := range tests {
		p := parse(tt.in)
		c := compiler.New()
		err := c.Compile(p)
		if err != nil {
			t.Fatalf("compiler got error: %s", err)
		}

		vm := New(c.ByteCode())
		err = vm.Run()

		if err != nil {
			t.Fatalf("vm.Run got error: %s", err)
		}

		stackElem := vm.LastPoppedStackElement()

		testExpectedObject(t, tt.want, stackElem)
	}
}

func testExpectedObject(
	t *testing.T,
	want interface{},
	got object.Object,
) {
	t.Helper()
	switch want := want.(type) {
	case int:
		err := testIntegerObject(int64(want), got)
		if err != nil {
			t.Errorf("testExpectedObject failed: %s", err)
		}

	}
}

func parse(in string) *ast.Program {
	l := lexer.New(in)
	p := parser.New(l)

	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	ret, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer.got=%T (%+v)", actual, actual)
	}
	if ret.Value != expected {
		return fmt.Errorf("object has wrong value. want=%d,got=%d", expected, ret.Value)
	}
	return nil
}
