package calc_test

import (
	Calc "calc"
	"math"
	"testing"
)

type Test struct {
	input    string
	expected float64
}

var tests = map[string]Test{
	"Test case 0": {
		input:    "( 1 + 2 ) / 3",
		expected: 1.0,
	},

	"Test case 1": {
		input:    "( 1 + 23 - 54 * 65 ) / 3 * 6",
		expected: -6972.0,
	},

	"Test case 2": {
		input:    "( 1 + 5 + 6 ) / ( 3 * 2 )",
		expected: 2.0,
	},

	"Test case 3": {
		input:    "( 77 + 99.245 - 1.0 ) / 5",
		expected: 35.049,
	},

	"Test case 4": {
		input:    "0 / 5.2",
		expected: 0.0,
	},

	"Test case 5": {
		input:    "5.2 / 0",
		expected: math.Inf(1),
	},

	"Test case 6": {
		input:    "-5.2 / 0",
		expected: math.Inf(-1),
	},
}

func TestCalc(t *testing.T) {
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			result, err := Calc.EvaluateRPN(Calc.ConvertToPolishReverseForm(test.input))
			if err != nil {
				t.Errorf("Calc returned an error: %v", err)
			}

			if result != test.expected {
				t.Errorf("Calc output does not match expected.\nExpected:\n%f\n\nGot:\n%f", test.expected, result)
			}
		})
	}
}
