package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Expression struct {
	value string
}

func parseExpression() Expression {
	expression := Expression{}
	flag.Parse()
	expression.value = flag.Arg(0)
	return expression
}

func evaluateRPN(expression string) (int, error) {
	stack := []int{}
	operators := map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}

	tokens := strings.Split(expression, " ")
	for _, token := range tokens {
		if operator, ok := operators[token]; ok {
			if len(stack) < 2 {
				return 0, fmt.Errorf("Invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			result := operator(a, b)
			stack = append(stack, result)
		} else {
			num, err := strconv.Atoi(token)
			if err != nil {
				return 0, fmt.Errorf("Invalid expression")
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("Invalid expression")
	}

	return stack[0], nil
}

func convertToPolishReverseForm(expression string) string {
	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	var result []string
	var stack []string

	tokens := strings.Split(expression, " ")
	for _, token := range tokens {
		if operator, ok := precedence[token]; ok {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if precedence[top] >= operator {
					result = append(result, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				result = append(result, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 && stack[len(stack)-1] == "(" {
				stack = stack[:len(stack)-1]
			}
		} else {
			result = append(result, token)
		}
	}

	for len(stack) > 0 {
		result = append(result, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return strings.Join(result, " ")
}

func main() {
	expression := parseExpression()
	if expression.value != "" {
		result, err := evaluateRPN(convertToPolishReverseForm(expression.value))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Result:", result)
		}
	} else {
		fmt.Println("Unknown expression")
	}
}
