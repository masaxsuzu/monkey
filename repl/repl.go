package repl

import (
	"bufio"
	"fmt"
	"github.com/masa-suzu/monkey/compiler"
	"github.com/masa-suzu/monkey/evaluator"
	"github.com/masa-suzu/monkey/lexer"
	"github.com/masa-suzu/monkey/object"
	"github.com/masa-suzu/monkey/parser"
	"github.com/masa-suzu/monkey/vm"
	"io"
)

func Start(in io.Reader, out io.Writer, prompt string, useVM bool) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if useVM {
			Rep_VM(line, out)
		} else {
			Rep(line, out, env)
		}
	}
}

func Rep(in string, out io.Writer, env *object.Environment) {
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrorsWithMonkeyFace(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func Rep_VM(in string, out io.Writer) {
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrorsWithMonkeyFace(out, p.Errors())
		return
	}

	c := compiler.New()

	err := c.Compile(program)

	if err != nil {
		printParserErrorsWithMonkeyFace(out, []string{err.Error()})
		return
	}

	vMachine := vm.New(c.ByteCode())
	err = vMachine.Run()

	if err != nil {
		printParserErrorsWithMonkeyFace(out, []string{err.Error()})
		return
	}

	lastPopped := vMachine.LastPoppedStackElement()
	if lastPopped != nil {
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrorsWithMonkeyFace(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
