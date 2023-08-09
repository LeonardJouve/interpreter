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
			case *object.Array:
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
	"first": {
		Value: func(arguments ...object.Object) object.Object {
			expectedArgumentAmount := 1
			if argumentAmout := len(arguments); argumentAmout != expectedArgumentAmount {
				return &object.Error{
					Value: fmt.Sprintf("wrong arguments amount: received %d, expected %d", argumentAmout, expectedArgumentAmount),
				}
			}

			if arguments[0].Type() != object.ARRAY {
				return &object.Error{
					Value: fmt.Sprintf("unsupported argument for builtin function first: %s", arguments[0].Type()),
				}
			}

			array, ok := arguments[0].(*object.Array)
			if !ok {
				return &object.Error{
					Value: fmt.Sprintf("invalid object type: received %T, expected *object.Array", arguments[0]),
				}
			}

			if len(array.Value) == 0 {
				return NULL
			}

			return array.Value[0]
		},
	},
	"last": {
		Value: func(arguments ...object.Object) object.Object {
			expectedArgumentAmount := 1
			if argumentAmout := len(arguments); argumentAmout != expectedArgumentAmount {
				return &object.Error{
					Value: fmt.Sprintf("wrong arguments amount: received %d, expected %d", argumentAmout, expectedArgumentAmount),
				}
			}

			if arguments[0].Type() != object.ARRAY {
				return &object.Error{
					Value: fmt.Sprintf("unsupported argument for builtin function last: %s", arguments[0].Type()),
				}
			}

			array, ok := arguments[0].(*object.Array)
			if !ok {
				return &object.Error{
					Value: fmt.Sprintf("invalid object type: received %T, expected *object.Array", arguments[0]),
				}
			}

			if len(array.Value) == 0 {
				return NULL
			}

			return array.Value[len(array.Value)-1]
		},
	},
	"rest": {
		Value: func(arguments ...object.Object) object.Object {
			expectedArgumentAmount := 1
			if argumentAmout := len(arguments); argumentAmout != expectedArgumentAmount {
				return &object.Error{
					Value: fmt.Sprintf("wrong arguments amount: received %d, expected %d", argumentAmout, expectedArgumentAmount),
				}
			}

			if arguments[0].Type() != object.ARRAY {
				return &object.Error{
					Value: fmt.Sprintf("unsupported argument for builtin function rest: %s", arguments[0].Type()),
				}
			}

			array, ok := arguments[0].(*object.Array)
			if !ok {
				return &object.Error{
					Value: fmt.Sprintf("invalid object type: received %T, expected *object.Array", arguments[0]),
				}
			}

			length := len(array.Value)
			if length == 0 {
				return NULL
			}

			elements := make([]object.Object, length-1)
			copy(elements, array.Value[1:length])

			return &object.Array{
				Value: elements,
			}
		},
	},
	"push": {
		Value: func(arguments ...object.Object) object.Object {
			expectedArgumentAmount := 2
			if argumentAmout := len(arguments); argumentAmout != expectedArgumentAmount {
				return &object.Error{
					Value: fmt.Sprintf("wrong arguments amount: received %d, expected %d", argumentAmout, expectedArgumentAmount),
				}
			}

			if arguments[0].Type() != object.ARRAY {
				return &object.Error{
					Value: fmt.Sprintf("unsupported argument for builtin function push: %s", arguments[0].Type()),
				}
			}

			array, ok := arguments[0].(*object.Array)
			if !ok {
				return &object.Error{
					Value: fmt.Sprintf("invalid object type: received %T, expected *object.Array", arguments[0]),
				}
			}
			length := len(array.Value)

			elements := make([]object.Object, length+1)
			copy(elements, array.Value)
			elements[length] = arguments[1]

			return &object.Array{
				Value: elements,
			}
		},
	},
}
