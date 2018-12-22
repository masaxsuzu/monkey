package vm

import (
	"fmt"
	"github.com/masa-suzu/monkey/code"
	"github.com/masa-suzu/monkey/compiler"
	"github.com/masa-suzu/monkey/object"
)

const StackSize = 2048

type VirtualMachine struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // Points to the next value. Top of stack is stack[sp-1]
}

func New(byteCode *compiler.ByteCode) *VirtualMachine {
	return &VirtualMachine{
		instructions: byteCode.Instructions,
		constants:    byteCode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VirtualMachine) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VirtualMachine) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			index := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[index])

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VirtualMachine) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}
