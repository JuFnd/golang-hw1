package uniq_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"uniq"
	Uniq "uniq"
)

type Test struct {
	input    string
	expected string
	flags    Uniq.UniqFlags
}

func parseTestFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return ""
	}

	str := string(data)
	return str
}

var tests = map[string]Test{
	"Test case 0": {
		input:    parseTestFile("tests/test.txt"),
		expected: "I love music.\n \nI love music of Kartik.\nThanKS.\nThanks.\nI love music of Kartik.\nabsdef\n",
		flags:    Uniq.UniqFlags{},
	},

	"Test case 1": {
		input:    parseTestFile("tests/test.txt"),
		expected: " \nThanKS.\nThanks.\nabsdef\n",
		flags: Uniq.UniqFlags{
			Unique: true,
		},
	},

	"Test case 2": {
		input:    parseTestFile("tests/test.txt"),
		expected: "I love music.\nI love music of Kartik.\nI love music of Kartik.\n",
		flags: Uniq.UniqFlags{
			Duplicates: true,
		},
	},

	"Test case 3": {
		input:    parseTestFile("tests/test.txt"),
		expected: "3 I love music.\n2 I love music of Kartik.\n2 I love music of Kartik.\n",
		flags: Uniq.UniqFlags{
			Duplicates: true,
			Count:      true,
		},
	},

	"Test case 4": {
		input:    parseTestFile("tests/test.txt"),
		expected: "1\n1 ThanKS.\n1 Thanks.\n1 absdef\n",
		flags: Uniq.UniqFlags{
			Count:  true,
			Unique: true,
		},
	},

	"Test case 5": {
		input:    parseTestFile("tests/test.txt"),
		expected: "3 I love music.\n1\n2 I love music of Kartik.\n1 ThanKS.\n1 Thanks.\n2 I love music of Kartik.\n1 absdef\n",
		flags: Uniq.UniqFlags{
			Count: true,
		},
	},

	"Test case 6": {
		input:    parseTestFile("tests/test.txt"),
		expected: "1\n1 absdef\n",
		flags: Uniq.UniqFlags{
			Count:      true,
			IgnoreCase: true,
			Unique:     true,
		},
	},

	"Test case 7": {
		input:    parseTestFile("tests/test.txt"),
		expected: "3 I love music.\n2 I love music of Kartik.\n2 Thanks.\n2 I love music of Kartik.\n",
		flags: Uniq.UniqFlags{
			Count:      true,
			IgnoreCase: true,
			Duplicates: true,
		},
	},

	"Test case 8": {
		input:    parseTestFile("tests/test.txt"),
		expected: "5\n2 I love music of Kartik.\n2 Thanks.\n2 I love music of Kartik.\n",
		flags: Uniq.UniqFlags{
			Count:      true,
			IgnoreCase: true,
			Duplicates: true,
			NumFields:  3,
		},
	},

	"Test case 9": {
		input:    parseTestFile("tests/test.txt"),
		expected: "3 I love music.\n2 I love music of Kartik.\n2 Thanks.\n2 I love music of Kartik.\n",
		flags: Uniq.UniqFlags{
			Count:      true,
			IgnoreCase: true,
			Duplicates: true,
			NumChars:   3,
		},
	},
}

func TestUniq(t *testing.T) {
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			input := strings.NewReader(test.input)
			output := bytes.NewBuffer(nil)

			err := uniq.Uniq(input, output, &test.flags)
			if err != nil {
				t.Errorf("Uniq returned an error: %v", err)
			}

			result := output.String()
			if result != test.expected {
				t.Errorf("Uniq output does not match expected.\nExpected:\n%s\n\nGot:\n%s", test.expected, result)
			}
		})
	}
}
