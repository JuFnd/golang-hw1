package calc

import (
	"flag"
	"fmt"
	"stack"
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
	stack := stack.NewStack()
	operators := map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}

	tokens := strings.Split(expression, " ")
	for _, token := range tokens {
		if operator, ok := operators[token]; ok {
			if stack.Size() < 2 {
				return 0, fmt.Errorf("Invalid expression")
			}
			b, _ := stack.Pop()
			a, _ := stack.Pop()
			result := operator(a.(int), b.(int))
			stack.Push(result)
		} else {
			num, err := strconv.Atoi(token)
			if err != nil {
				return 0, fmt.Errorf("Invalid expression")
			}
			stack.Push(num)
		}
	}

	if stack.Size() != 1 {
		return 0, fmt.Errorf("Invalid expression")
	}

	result, _ := stack.Pop()
	return result.(int), nil
}

func convertToPolishReverseForm(expression string) string {
	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	var result []string
	var stack stack.Stack

	tokens := strings.Split(expression, " ")
	for _, token := range tokens {
		if operator, ok := precedence[token]; ok {
			for stack.Size() > 0 {
				top, _ := stack.Pop()
				if precedence[top.(string)] >= operator {
					result = append(result, top.(string))
				} else {
					stack.Push(top)
					break
				}
			}
			stack.Push(token)
		} else if token == "(" {
			stack.Push(token)
		} else if token == ")" {
			for stack.Size() > 0 {
				top, _ := stack.Pop()
				if top.(string) == "(" {
					break
				}
				result = append(result, top.(string))
			}
		} else {
			result = append(result, token)
		}
	}

	for stack.Size() > 0 {
		top, _ := stack.Pop()
		result = append(result, top.(string))
	}

	return strings.Join(result, " ")
}

func Calc() {
	expression := parseExpression()
	if expression.value == "" {
		fmt.Println("Unknown expression")
		return
	}

	result, err := evaluateRPN(convertToPolishReverseForm(expression.value))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Result:", result)
}
