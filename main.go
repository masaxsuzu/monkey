package main

import (
	"flag"
	"fmt"
	"github.com/masa-suzu/monkey/repl"
	"os"
	"os/user"
)

var (
	useVM = flag.Bool("vm", false, "run on virtual machine")
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Execute exit() then exit monkey!\n")

	repl.Start(os.Stdin, os.Stdout, "$ ", *useVM)
}
