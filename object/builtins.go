package object

import (
	"fmt"
	"os"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		Name: "len",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				default:
					return newError("argument to `len` not supported, got %s", arg.Type())
				}
			},
		},
	},
	{
		Name: "puts",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
	},
	{
		Name: "first",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of argument. got=%d, want=1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to first must be ARRAY, got %s", args[0].Type())
				}
				array := args[0].(*Array)
				length := len(array.Elements)
				if length > 0 {
					return array.Elements[0]
				}
				return nil
			},
		},
	},
	{
		Name: "last",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of argument. got=%d, want=1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to last must be ARRAY, got %s", args[0].Type())
				}
				array := args[0].(*Array)
				length := len(array.Elements)
				if length > 0 {
					return array.Elements[length-1]
				}
				return nil
			},
		},
	},
	{
		Name: "rest",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of argument. got=%d, want=1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to rest must be ARRAY, got %s", args[0].Type())
				}
				array := args[0].(*Array)
				length := len(array.Elements)
				if length > 0 {
					newElements := make([]Object, length-1, length-1)
					copy(newElements, array.Elements[1:length])
					return &Array{Elements: newElements}
				}
				return nil
			},
		},
	},
	{
		Name: "help",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				fmt.Println("This is the Monkey programming language!")
				fmt.Println("Execute exit() then exit monkey!")
				return nil
			},
		},
	},
	{
		Name: "exit",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				os.Exit(0)
				return nil
			},
		},
	},
}

func GetBuiltinName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
