package benchmark

import (
	"fmt"
	"github.com/masa-suzu/monkey/ast"
	"github.com/masa-suzu/monkey/compiler"
	"github.com/masa-suzu/monkey/evaluator"
	"github.com/masa-suzu/monkey/lexer"
	"github.com/masa-suzu/monkey/object"
	"github.com/masa-suzu/monkey/parser"
	"github.com/masa-suzu/monkey/vm"
	"testing"
)

var input string = `
let f = fn(x){
    if (x < 2) { return x}
    return f(x-1) + f(x-2)
}
f(%v)
`

func BenchmarkRun10NestedFibonacci_Evaluator(b *testing.B) {
	evaluate(fmt.Sprintf(input, 10))
}
func BenchmarkRun10NestedFibonacci_VM(b *testing.B) {
	run(fmt.Sprintf(input, 10))
}

func BenchmarkRun20NestedFibonacci_Evaluator(b *testing.B) {
	evaluate(fmt.Sprintf(input, 20))
}
func BenchmarkRun20NestedFibonacci_VM(b *testing.B) {
	run(fmt.Sprintf(input, 20))
}

func BenchmarkRun25NestedFibonacci_Evaluator(b *testing.B) {
	evaluate(fmt.Sprintf(input, 25))
}
func BenchmarkRun25NestedFibonacci_VM(b *testing.B) {
	run(fmt.Sprintf(input, 25))
}

func BenchmarkRun26NestedFibonacci_Evaluator(b *testing.B) {
	evaluate(fmt.Sprintf(input, 26))
}
func BenchmarkRun26NestedFibonacci_VM(b *testing.B) {
	run(fmt.Sprintf(input, 26))
}

func BenchmarkRun27NestedFibonacci_Evaluator(b *testing.B) {
	evaluate(fmt.Sprintf(input, 27))
}
func BenchmarkRun27NestedFibonacci_VM(b *testing.B) {
	run(fmt.Sprintf(input, 27))
}

func BenchmarkRun28NestedFibonacci_Evaluator(b *testing.B) {
	evaluate(fmt.Sprintf(input, 28))
}
func BenchmarkRun28NestedFibonacci_VM(b *testing.B) {
	run(fmt.Sprintf(input, 28))
}

func BenchmarkRun29NestedFibonacci_Evaluator(b *testing.B) {
	evaluate(fmt.Sprintf(input, 29))
}
func BenchmarkRun29NestedFibonacci_VM(b *testing.B) {
	run(fmt.Sprintf(input, 29))
}

func BenchmarkRun30NestedFibonacci_Evaluator(b *testing.B) {
	evaluate(fmt.Sprintf(input, 30))
}
func BenchmarkRun30NestedFibonacci_VM(b *testing.B) {
	run(fmt.Sprintf(input, 30))
}

func evaluate(src string) {
	p := parse(src)
	env := object.NewEnvironment()
	evaluator.Eval(p, env)
}

func run(src string) {
	p := parse(src)

	c := compiler.New()
	c.Compile(p)

	vm := vm.New(c.ByteCode())

	vm.Run()
	vm.LastPoppedStackElement()
}

func parse(src string) *ast.Program {
	l := lexer.New(src)
	p := parser.New(l)
	return p.ParseProgram()
}
