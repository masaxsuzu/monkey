// +build js

package main

import (
	"bytes"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/masa-suzu/monkey/repl"
)

func main() {
	js.Global.Set("run", startRepl)
}

func startRepl(source string) string {
	in := bytes.NewBufferString(source)
	out := bytes.NewBufferString("")
	repl.Start(in, out, "")
	return fmt.Sprint(out)
}
