package main

import (
	"bytes"
	"fmt"
	"github.com/masa-suzu/monkey/compiler"
	"github.com/masa-suzu/monkey/formatter"
	"github.com/masa-suzu/monkey/lexer"
	"github.com/masa-suzu/monkey/parser"
	"github.com/masa-suzu/monkey/vm"

	"github.com/gopherjs/gopherjs/js"
	"github.com/masa-suzu/monkey/object"
	"github.com/masa-suzu/monkey/repl"
)

func main() {
	js.Global.Set("run", startRep)
	js.Global.Set("fmt", startFormat)

}

func startRep(source string) string {
	out := bytes.NewBufferString("")
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}
	repl.Rep_VM(source, out, false, constants, globals, symbolTable)
	return fmt.Sprint(out)
}

func startFormat(source string) string {
	out := bytes.NewBufferString("")
	l := lexer.New(source)
	p := parser.New(l)
	ast := p.ParseProgram()
	out.WriteString(formatter.Format(ast, 0))
	return fmt.Sprint(out)
}
