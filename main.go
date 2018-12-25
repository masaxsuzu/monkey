package main

import (
	"flag"
	"fmt"
	"github.com/masa-suzu/monkey/repl"
	"os"
)

var (
	useVM     = flag.Bool("vm", false, "run on virtual machine")
	debugMode = flag.Bool("debug", false, "dump instructions on virtual machine for each run")
)

func main() {

	flag.Parse()

	fmt.Printf("Hello! This is the Monkey programming language!\n")

	repl.Start(os.Stdin, os.Stdout, ">> ", *useVM, *debugMode)
}
