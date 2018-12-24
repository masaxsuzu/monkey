package evaluator

import (
	"fmt"
	"github.com/masa-suzu/monkey/ast"
	"github.com/masa-suzu/monkey/object"
	"github.com/masa-suzu/monkey/token"
)

func quote(node ast.Node, env *object.Environment) *object.Quote {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}
func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return nodeFrom(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExp, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExp.Function.TokenLiteral() == "unquote"
}

func nodeFrom(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}

	default:
		return nil
	}
}
