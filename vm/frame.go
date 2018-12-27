package vm

import (
	"github.com/masa-suzu/monkey/code"
	"github.com/masa-suzu/monkey/object"
)

type Frame struct {
	closure     *object.Closure
	ip          int
	basePointer int
}

func NewFrame(c *object.Closure, basePointer int) *Frame {
	return &Frame{
		closure:     c,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.closure.Function.Instructions
}
