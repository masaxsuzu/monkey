package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type OperandCode byte

type Definition struct {
	Name          string
	OperandWidths []int
}

const (
	Constant = iota
	Add
	Sub
	Mul
	Div
	Pop
	True
	False
	Equal
	NotEqual
	GreaterThan
	Minus
	Bang
	Jump
	JumpNotTruthy
)

var definitions = map[OperandCode]*Definition{
	Constant:      {"Constant", []int{2}},
	Add:           {"Add", []int{}},
	Sub:           {"Sub", []int{}},
	Mul:           {"Mul", []int{}},
	Div:           {"Div", []int{}},
	Pop:           {"Pop", []int{}},
	True:          {"True", []int{}},
	False:         {"False", []int{}},
	Equal:         {"Equal", []int{}},
	NotEqual:      {"NotEqual", []int{}},
	GreaterThan:   {"GreaterThan", []int{}},
	Minus:         {"Minus", []int{}},
	Bang:          {"Bang", []int{}},
	Jump:          {"Jump", []int{2}},
	JumpNotTruthy: {"JumpNotTruthy", []int{2}},
}

func (ins Instructions) String() string {
	var out bytes.Buffer
	i := 0

	for i < len(ins) {
		def, err := LookUp(ins[i])

		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func LookUp(op byte) (*Definition, error) {
	def, ok := definitions[(OperandCode(op))]

	if !ok {
		return nil, fmt.Errorf("opcode %d undifined", op)
	}
	return def, nil
}

func Make(op OperandCode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))

	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
