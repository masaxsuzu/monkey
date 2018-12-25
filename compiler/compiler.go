package compiler

import (
	"fmt"
	"github.com/masa-suzu/monkey/ast"
	"github.com/masa-suzu/monkey/code"
	"github.com/masa-suzu/monkey/object"
	"sort"
)

type EmittedInstruction struct {
	Code     code.OperandCode
	Position int
}

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
}

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
	}
}

func NewWithState(st *SymbolTable, constants []object.Object) *Compiler {
	c := New()
	c.symbolTable = st
	c.constants = constants
	return c
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
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.SetGlobal, symbol.Index)
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
	case *ast.IfExpression:
		var err error = nil
		compileNode := func(n ast.Node) {
			if err == nil {
				err = c.Compile(n)
			}
		}

		compileNode(node.Condition)

		if err != nil {
			return err
		}

		jumpNotTruthyPos := c.emit(code.JumpNotTruthy, -1)

		compileNode(node.Consequence)

		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		changeOperandAtXByTail := func(x int) {
			tailPos := len(c.instructions)
			c.changeOperand(x, tailPos)
		}

		jumpPos := c.emit(code.Jump, -1)

		changeOperandAtXByTail(jumpNotTruthyPos)

		if node.Alternative == nil {
			c.emit(code.Null)
		} else {

			err = c.Compile(node.Alternative)

			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}

		}

		changeOperandAtXByTail(jumpPos)
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.emit(code.GetGlobal, symbol.Index)
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.Constant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.True)
		} else {
			c.emit(code.False)
		}
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.Constant, c.addConstant(str))
	case *ast.ArrayLiteral:
		for _, v := range node.Elements {
			err := c.Compile(v)

			if err != nil {
				return err
			}
		}
		c.emit(code.Array, len(node.Elements))
	case *ast.HashLiteral:
		var err error = nil
		compileNode := func(n ast.Node) {
			if err == nil {
				err = c.Compile(n)
			}
		}
		keys := []ast.Expression{}
		for key := range node.Pairs {
			keys = append(keys, key)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, key := range keys {
			compileNode(key)
			compileNode(node.Pairs[key])
			if err != nil {
				return err
			}
		}

		c.emit(code.Hash, len(node.Pairs)*2)
	case *ast.IndexExpression:
		var err error = nil
		compileNode := func(n ast.Node) {
			if err == nil {
				err = c.Compile(n)
			}
		}
		compileNode(node.Left)
		compileNode(node.Index)

		if err != nil {
			return err
		}
		c.emit(code.Index)
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

	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op code.OperandCode, pos int) {
	prev := c.lastInstruction
	last := EmittedInstruction{op, pos}

	c.previousInstruction = prev
	c.lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Code == code.Pop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(opPos int, new []byte) {
	for i := 0; i < len(new); i++ {
		c.instructions[opPos+i] = new[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.OperandCode(c.instructions[opPos])
	newIns := code.Make(op, operand)
	c.replaceInstruction(opPos, newIns)
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}
