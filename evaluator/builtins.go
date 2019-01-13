package evaluator

import (
	"github.com/masa-suzu/monkey/object"
)

var builtIns = map[string]*object.Builtin{
	"len":   object.GetBuiltinName("len"),
	"last":  object.GetBuiltinName("last"),
	"first": object.GetBuiltinName("first"),
	"rest":  object.GetBuiltinName("rest"),
	"puts":  object.GetBuiltinName("puts"),
	"exit":  object.GetBuiltinName("exit"),
	"help":  object.GetBuiltinName("help"),
}
