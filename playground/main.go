package main

import (
	"bytes"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/masa-suzu/monkey/repl"
	"github.com/masa-suzu/monkey/object"
)

func main() {
	js.Global.Set("run", startRep)
}

func startRep(source string) string {
	out := bytes.NewBufferString("")
	env := object.NewEnvironment()
	repl.Rep(source,out, "", env)
	return fmt.Sprint(out)
}
