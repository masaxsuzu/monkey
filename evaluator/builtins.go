package evaluator

import (
	"fmt"
	"monkey/object"
	"os"
)

var builtIns = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of argument. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to last must be ARRAY, got %s", args[0].Type())
			}
			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length > 0 {
				return array.Elements[length-1]
			}
			return NULL
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of argument. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}
			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length > 0 {
				return array.Elements[0]
			}
			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of argument. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to rest must be ARRAY, got %s", args[0].Type())
			}
			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, array.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"exit": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			os.Exit(0)
			return NULL
		},
	},
	"help": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			fmt.Println("This is the Monkey programming language!")
			fmt.Println("Execute exit() then exit monkey!")
			return NULL
		},
	},
}
