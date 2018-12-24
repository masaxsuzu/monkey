package main

import (
	"bytes"
	"fmt"
	"github.com/masa-suzu/monkey/formatter"
	"github.com/masa-suzu/monkey/lexer"
	"github.com/masa-suzu/monkey/parser"

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
	env := object.NewEnvironment()
	macros := object.NewEnvironment()
	repl.Rep(source, out, env,macros)
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
