package main

import (
	"flag"
	"fmt"
	"github.com/masa-suzu/monkey/repl"
	"os"
	"os/user"
)

var (
	useVM     = flag.Bool("vm", false, "run on virtual machine")
	debugMode = flag.Bool("debug", false, "dump instructions on virtual machine for each run")
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", u.Username)
	fmt.Printf("Execute exit() then exit monkey!\n")

	repl.Start(os.Stdin, os.Stdout, "$ ", *useVM, *debugMode)
}
