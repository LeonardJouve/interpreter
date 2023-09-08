package ast

type Modifier func(Node) Node

func Modify(node Node, modifier Modifier) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *BlockStatement:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *ReturnStatement:
		node.Value = Modify(node.Value, modifier).(Expression)
	case *LetStatement:
		node.Value = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i, element := range node.Value {
			node.Value[i] = Modify(element, modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, value := range node.Value {
			newKey, _ := Modify(key, modifier).(Expression)
			newValue, _ := Modify(value, modifier).(Expression)
			newPairs[newKey] = newValue
		}
		node.Value = newPairs
	}

	return modifier(node)
}
