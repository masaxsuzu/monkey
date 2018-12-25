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

func Start(in io.Reader, out io.Writer, prompt string, useVM bool, debugMode bool) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()
	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if useVM {
			Rep_VM(line, out, debugMode, constants, globals, symbolTable)
		} else {
			Rep(line, out, env, macroEnv)
		}
	}
}

func Rep(in string, out io.Writer, env *object.Environment, macros *object.Environment) {
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printErrorsWithMonkeyFace(out, p.Errors(), "Parser")
		return
	}

	evaluator.DefineMacros(program, macros)
	expanded := evaluator.ExpandMacros(program, macros)

	evaluated := evaluator.Eval(expanded, env)

	if evaluated != nil {
		if evaluated.Type() == object.ERROR_OBJ {
			printErrorsWithMonkeyFace(out, []string{evaluated.Inspect()}, "Run time")
			return
		}
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func Rep_VM(in string, out io.Writer, debugMode bool, constants []object.Object, scope []object.Object, st *compiler.SymbolTable) {
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printErrorsWithMonkeyFace(out, p.Errors(), "Parser")
		return
	}

	c := compiler.NewWithState(st, constants)

	err := c.Compile(program)

	if err != nil {
		printErrorsWithMonkeyFace(out, []string{err.Error()}, "Compile")
		return
	}

	var vMachine *vm.VirtualMachine
	code := c.ByteCode()
	constants = code.Constants
	vMachine = vm.NewWithGlobalScope(code, scope)
	vMachine.DebugMode = debugMode
	err = vMachine.Run()

	if err != nil {
		printErrorsWithMonkeyFace(out, []string{err.Error()}, "Run time")
		return
	}

	lastPopped := vMachine.LastPoppedStackElement()
	if lastPopped != nil {
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
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

func printErrorsWithMonkeyFace(out io.Writer, errors []string, label string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, fmt.Sprintf("%s errors:\n", label))
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
