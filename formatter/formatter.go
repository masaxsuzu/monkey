package formatter

import (
	"bytes"
	"strings"

	"github.com/masa-suzu/monkey/ast"
)

func Format(node ast.Node, indent int) string {
	switch v := node.(type) {
	case *ast.Program:
		out := bytes.Buffer{}
		for i, s := range v.Statements {
			if i != 0 {
				out.WriteString("\n")
			}
			out.WriteString(Format(s, indent))
		}
		return out.String()
	case *ast.ExpressionStatement:
		return Format(v.Expression, indent) + ";"
	case *ast.BlockStatement:
		out := bytes.Buffer{}
		for i, s := range v.Statements {
			if i != 0 {
				out.WriteString("\n")
			}
			out.WriteString(Format(s, indent))
		}
		return out.String()
	case *ast.IntegerLiteral:
		return indents(indent) + v.String()
	case *ast.StringLiteral:
		return indents(indent) + v.String()
	case *ast.Boolean:
		return indents(indent) + v.String()
	case *ast.PrefixExpression:
		return indents(indent) + v.String()
	case *ast.InfixExpression:
		return indents(indent) + v.String()
	case *ast.IfExpression:
		out := bytes.Buffer{}
		out.WriteString(indents(indent) + "if(")
		out.WriteString(Format(v.Condition, 0))
		out.WriteString(") {\n")
		out.WriteString(Format(v.Consequence, indent+1))
		out.WriteString("\n")
		out.WriteString(indents(indent) + "}")
		if v.Alternative != nil {
			out.WriteString(indents(indent) + " else {\n")
			out.WriteString(Format(v.Alternative, indent+1))
			out.WriteString("\n")
			out.WriteString(indents(indent) + "}")
		}
		return out.String()
	case *ast.ReturnStatement:
		out := bytes.Buffer{}
		out.WriteString(indents(indent) + "return ")
		out.WriteString(strings.Replace(Format(v.ReturnValue, indent)+";", indents(indent), "", 1))
		return out.String()
	case *ast.LetStatement:
		out := bytes.Buffer{}
		out.WriteString(indents(indent) + "let ")
		out.WriteString(v.Name.String())
		out.WriteString(" = ")
		out.WriteString(strings.Replace(Format(v.Value, indent)+";", indents(indent), "", 1))
		return out.String()
	case *ast.FunctionLiteral:
		var out bytes.Buffer

		params := []string{}

		for _, p := range v.Parameters {
			params = append(params, p.String())
		}
		out.WriteString(indents(indent) + "fn(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(") {\n")
		x := Format(v.Body, indent+1)
		out.WriteString(x)
		out.WriteString("\n")
		out.WriteString(indents(indent) + "}")
		return out.String()
	case *ast.MacroLiteral:
		var out bytes.Buffer

		params := []string{}

		for _, p := range v.Parameters {
			params = append(params, p.String())
		}
		out.WriteString(indents(indent) + "macro(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(") {\n")
		x := Format(v.Body, indent+1)
		out.WriteString(x)
		out.WriteString("\n")
		out.WriteString(indents(indent) + "}")
		return out.String()
	case *ast.CallExpression:
		// TODO: for quote/unquote
		return indents(indent) + v.String()
	case *ast.ArrayLiteral:
		return indents(indent) + v.String()
	case *ast.IndexExpression:
		return indents(indent) + v.String()
	case *ast.HashLiteral:
		return indents(indent) + v.String()
	case *ast.Identifier:
		return indents(indent) + v.String()
	default:
		return indents(indent) + v.String()
	}
}

func indents(level int) string {
	spaces := "    "
	out := &bytes.Buffer{}
	for i := 0; i < level; i++ {
		out.WriteString(spaces)
	}
	return out.String()
}
