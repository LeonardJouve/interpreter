package object

import "testing"

func TestStringHashKey(t *testing.T) {
	tests := [][]*String{
		{
			&String{
				Value: "hello world",
			},
			&String{
				Value: "hello world",
			},
		},
		{
			&String{
				Value: "test test test",
			},
			&String{
				Value: "test test test",
			},
		},
	}

	for _, test := range tests {
		expectedElementAmount := 2
		if elementAmount := len(test); elementAmount != expectedElementAmount {
			t.Errorf("[Test] Invalid element amount: received %d, expected %d", elementAmount, expectedElementAmount)
			continue
		}

		expectedHash := test[1].HashKey()
		if hash := test[0].HashKey(); hash != expectedHash {
			t.Errorf("[Test] Invalid string hash key: received %d, expected %d", hash.Value, expectedHash.Value)
		}
	}
}
