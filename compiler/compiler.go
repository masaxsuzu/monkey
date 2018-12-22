package compiler

import (
	"github.com/masa-suzu/monkey/ast"
	"github.com/masa-suzu/monkey/code"
	"github.com/masa-suzu/monkey/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
