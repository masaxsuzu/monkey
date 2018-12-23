package vm

import (
	"fmt"
	"github.com/masa-suzu/monkey/code"
	"github.com/masa-suzu/monkey/compiler"
	"github.com/masa-suzu/monkey/object"
)

const StackSize = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

type VirtualMachine struct {
	constants    []object.Object
	instructions code.Instructions

	stack     []object.Object
	sp        int // Points to the next value. Top of stack is stack[sp-1]
	debugMode bool
}

func New(byteCode *compiler.ByteCode) *VirtualMachine {
	return &VirtualMachine{
		instructions: byteCode.Instructions,
		constants:    byteCode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
		debugMode:    false,
	}
}

func DebugMode(byteCode *compiler.ByteCode) *VirtualMachine {
	vm := New(byteCode)
	vm.debugMode = true
	return vm
}

func (vm *VirtualMachine) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VirtualMachine) Run() error {
	if vm.debugMode {
		vm.dumpInstructions()
	}
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.OperandCode(vm.instructions[ip])
		switch op {
		case code.Constant:
			index := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[index])

			if err != nil {
				return err
			}
		case code.Add, code.Sub, code.Mul, code.Div:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.Equal, code.NotEqual, code.GreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case code.Pop:
			vm.pop()
		case code.True:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.False:
			err := vm.push(False)
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

func (vm *VirtualMachine) pop() object.Object {
	top := vm.stack[vm.sp-1]
	vm.sp--
	return top
}

func (vm *VirtualMachine) LastPoppedStackElement() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VirtualMachine) executeBinaryOperation(op code.OperandCode) error {
	r := vm.pop()
	l := vm.pop()

	rightType := r.Type()
	leftType := l.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, l.(*object.Integer), r.(*object.Integer))
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VirtualMachine) executeBinaryIntegerOperation(op code.OperandCode, left *object.Integer, right *object.Integer) error {
	lv := left.Value
	rv := right.Value
	var ret int64
	switch op {
	case code.Add:
		ret = lv + rv
	case code.Sub:
		ret = lv - rv
	case code.Mul:
		ret = lv * rv
	case code.Div:
		if rv == 0 {
			return fmt.Errorf("integer divide by zero")
		}
		ret = lv / rv
	default:
		return fmt.Errorf("uknown integer operator: %d", op)
	}
	return vm.push(&object.Integer{Value: ret})
}

func (vm *VirtualMachine) executeComparison(op code.OperandCode) error {
	r := vm.pop()
	l := vm.pop()

	if l.Type() == object.INTEGER_OBJ || r.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, l, r)
	}

	switch op {
	case code.Equal:
		return vm.push(nativeBoolToBooleanObject(r == l))
	case code.NotEqual:
		return vm.push(nativeBoolToBooleanObject(r != l))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, l.Type(), r.Type())
	}
}

func (vm *VirtualMachine) executeIntegerComparison(
	op code.OperandCode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.Equal:
		return vm.push(nativeBoolToBooleanObject(leftValue == rightValue))
	case code.NotEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue != rightValue))
	case code.GreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBoolToBooleanObject(ret bool) *object.Boolean {
	if ret {
		return True
	}
	return False
}

func (vm *VirtualMachine) dumpInstructions() {
	fmt.Printf("[dump]\n%s[out]\n", vm.instructions)
}
