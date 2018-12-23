package compiler

import (
	"fmt"
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
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.Pop)
	case *ast.InfixExpression:
		var err error = nil
		compileNode := func(n ast.Node) {
			if err == nil {
				err = c.Compile(n)
			}
		}
		if node.Operator == "<" {
			compileNode(node.Right)
			compileNode(node.Left)
		} else {
			compileNode(node.Left)
			compileNode(node.Right)
		}
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.Add)
		case "-":
			c.emit(code.Sub)
		case "*":
			c.emit(code.Mul)
		case "/":
			c.emit(code.Div)
		case "==":
			c.emit(code.Equal)
		case "!=":
			c.emit(code.NotEqual)
		case ">":
			c.emit(code.GreaterThan)
		case "<":
			c.emit(code.GreaterThan)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.PrefixExpression:
		var err error = nil
		compileNode := func(n ast.Node) {
			if err == nil {
				err = c.Compile(n)
			}
		}

		compileNode(node.Right)

		if err != nil {
			return err
		}

		switch node.Operator {
		case "-":
			c.emit(code.Minus)
		case "!":
			c.emit(code.Bang)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.Constant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.True)
		} else {
			c.emit(code.False)
		}
	}
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

func (c *Compiler) emit(op code.OperandCode, operand ...int) int {
	ins := code.Make(op, operand...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}
