package evaluator

import (
	"fmt"
	"leonardjouve/object"
	"leonardjouve/token"
)

var builtins = map[token.TokenLiteral]*object.Builtin{
	"len": {
		Value: func(arguments ...object.Object) object.Object {
			expectedArgumentAmount := 1
			if argumentAmout := len(arguments); argumentAmout != expectedArgumentAmount {
				return &object.Error{
					Value: fmt.Sprintf("wrong arguments amount: received %d, expected %d", argumentAmout, expectedArgumentAmount),
				}
			}

			switch argument := arguments[0].(type) {
			case *object.String:
				return &object.Integer{
					Value: int64(len(argument.Value)),
				}
			default:
				return &object.Error{
					Value: fmt.Sprintf("unsupported argument for builtin function len: %s", argument.Type()),
				}
			}
		},
	},
}
