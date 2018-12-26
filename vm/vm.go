package vm

import (
	"fmt"
	"github.com/masa-suzu/monkey/code"
	"github.com/masa-suzu/monkey/compiler"
	"github.com/masa-suzu/monkey/object"
)

const StackSize = 2048
const GlobalSize = 65536

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VirtualMachine struct {
	DebugMode    bool
	constants    []object.Object
	instructions code.Instructions

	stack   []object.Object
	sp      int // Points to the next value. Top of stack is stack[sp-1]
	globals []object.Object
}

func New(byteCode *compiler.ByteCode) *VirtualMachine {
	return &VirtualMachine{
		instructions: byteCode.Instructions,
		constants:    byteCode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
		globals:      make([]object.Object, GlobalSize),
		DebugMode:    false,
	}
}

func NewWithGlobalScope(byteCode *compiler.ByteCode, s []object.Object) *VirtualMachine {
	vm := New(byteCode)
	vm.globals = s
	return vm
}

func (vm *VirtualMachine) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VirtualMachine) Run() error {
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
		case code.Minus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}
		case code.Bang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.Index:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}
		case code.JumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}

		case code.Jump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case code.GetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.SetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.Array:
			numElements := int(code.ReadUint16(vm.instructions[ip+1:]))
			if numElements == 0 {

			}
			ip += 2
			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)

			if err != nil {
				return err
			}
		case code.Hash:
			numElements := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElements

			err = vm.push(hash)

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
		case code.Null:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		}
	}
	if vm.DebugMode {
		vm.dump()
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
	if leftType == object.STRING_OBJ && rightType == object.STRING_OBJ {
		return vm.executeBinaryStringOperation(op, l.(*object.String), r.(*object.String))
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

func (vm *VirtualMachine) executeBinaryStringOperation(op code.OperandCode, left *object.String, right *object.String) error {
	lv := left.Value
	rv := right.Value
	var ret string
	switch op {
	case code.Add:
		ret = lv + rv
	default:
		return fmt.Errorf("uknown string operator: %d", op)
	}
	return vm.push(&object.String{Value: ret})
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
	l, ok := left.(*object.Integer)
	if !ok {
		return vm.push(nativeBoolToBooleanObject(op != code.Equal))
	}
	r, ok := right.(*object.Integer)
	if !ok {
		return vm.push(nativeBoolToBooleanObject(op != code.Equal))
	}
	leftValue := l.Value
	rightValue := r.Value
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

func (vm *VirtualMachine) executeMinusOperator() error {
	op := vm.pop()

	if op.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", op.Type())
	}
	value := op.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VirtualMachine) executeBangOperator() error {
	op := vm.pop()

	switch op {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VirtualMachine) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}
func (vm *VirtualMachine) executeArrayIndex(array, index object.Object) error {
	arrayObj := array.(*object.Array)
	i := index.(*object.Integer).Value

	max := int64(len(arrayObj.Elements) - 1)
	if i < 0 || max < i {
		return vm.push(Null)
	}
	return vm.push(arrayObj.Elements[i])
}

func (vm *VirtualMachine) executeHashIndex(hash, index object.Object) error {
	hashObj := hash.(*object.Hash)
	key, ok := index.(object.Hashable)

	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]

	if !ok {
		return vm.push(Null)
	}
	return vm.push(pair.Value)
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(ret bool) *object.Boolean {
	if ret {
		return True
	}
	return False
}

func (vm *VirtualMachine) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)
	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}
	return &object.Array{Elements: elements}
}

func (vm *VirtualMachine) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)

		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}
		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VirtualMachine) dump() {
	fmt.Printf("[instructions]\n%s", vm.instructions)
	fmt.Println("[global scope]")
	for i, v := range vm.globals {
		if v != nil {
			fmt.Printf("%04d %v\n", i, v.Inspect())
		}
	}
	fmt.Println("[out]")
}
