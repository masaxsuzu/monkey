package ast

type NodeModifier func(Node) Node

func Modify(from Node, modify NodeModifier) Node {
	// TODO: remove _ in type assertions.
	switch from := from.(type) {
	case *Program:
		for i, _ := range from.Statements {
			from.Statements[i], _ = Modify(from.Statements[i], modify).(Statement)
		}
	case *BlockStatement:
		for i, statement := range from.Statements {
			from.Statements[i], _ = Modify(statement, modify).(Statement)
		}
	case *ReturnStatement:
		from.ReturnValue, _ = Modify(from.ReturnValue, modify).(Expression)
	case *LetStatement:
		from.Value, _ = Modify(from.Value, modify).(Expression)
	case *ExpressionStatement:
		from.Expression, _ = Modify(from.Expression, modify).(Expression)
	case *InfixExpression:
		from.Left, _ = Modify(from.Left, modify).(Expression)
		from.Right, _ = Modify(from.Right, modify).(Expression)
	case *PrefixExpression:
		from.Right, _ = Modify(from.Right, modify).(Expression)
	case *IndexExpression:
		from.Left, _ = Modify(from.Left, modify).(Expression)
		from.Index, _ = Modify(from.Index, modify).(Expression)
	case *IfExpression:
		from.Condition, _ = Modify(from.Condition, modify).(Expression)
		from.Consequence, _ = Modify(from.Consequence, modify).(*BlockStatement)
		if from.Alternative != nil {
			from.Alternative, _ = Modify(from.Alternative, modify).(*BlockStatement)
		}
	case *FunctionLiteral:
		for i, _ := range from.Parameters {
			from.Parameters[i], _ = Modify(from.Parameters[i], modify).(*Identifier)
		}
		from.Body, _ = Modify(from.Body, modify).(*BlockStatement)
	case *ArrayLiteral:
		for i, _ := range from.Elements {
			from.Elements[i], _ = Modify(from.Elements[i], modify).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, val := range from.Pairs {
			newKey, _ := Modify(key, modify).(Expression)
			newVal, _ := Modify(val, modify).(Expression)
			newPairs[newKey] = newVal
		}
		from.Pairs = newPairs
	}

	return modify(from)
}
