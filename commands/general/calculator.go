package general

import (
	"errors"
	"github.com/ebarkie/aprs"
	"math"
	"simpleAPRSbot-go/helpers/aprsHelper"
	"strconv"
	"strings"
	"unicode"
)

func CalculateCommand(args []string, f aprs.Frame, client aprsHelper.APRSUserClient) {
	// this is a calculator. takes in a string, and returns an answer
	var input = strings.Join(args, " ")
	calculate, err := Calculate(input)
	if err != nil {
		client.AprsTextReply(err.Error(), f)
		return
	} else {
		client.AprsTextReply(strconv.FormatFloat(calculate, 'g', 5, 64), f)
		return
	}
}

// Token represents a single element in the expression
type Token struct {
	typ   string
	value string
}

// precedence defines the precedence of operators
var precedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"^": 3,
	"(": 0,
	")": 0,
}

// Calculate evaluates a mathematical expression string following PEMDAS rules
func Calculate(input string) (float64, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return 0, err
	}

	postfix, err := toPostfix(tokens)
	if err != nil {
		return 0, err
	}

	result, err := evaluatePostfix(postfix)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// tokenize splits the input string into tokens
func tokenize(input string) ([]Token, error) {
	var tokens []Token
	var number strings.Builder

	for _, char := range input {
		if unicode.IsDigit(char) || char == '.' {
			number.WriteRune(char)
		} else if char == ' ' {
			continue
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, Token{typ: "number", value: number.String()})
				number.Reset()
			}
			if strings.ContainsRune("+-*/^()", char) {
				tokens = append(tokens, Token{typ: "operator", value: string(char)})
			} else {
				return nil, errors.New("invalid character in input")
			}
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, Token{typ: "number", value: number.String()})
	}

	return tokens, nil
}

// toPostfix converts an infix expression to postfix notation
func toPostfix(tokens []Token) ([]Token, error) {
	var output []Token
	var stack []Token

	for _, token := range tokens {
		switch token.typ {
		case "number":
			output = append(output, token)
		case "operator":
			if token.value == "(" {
				stack = append(stack, token)
			} else if token.value == ")" {
				for len(stack) > 0 && stack[len(stack)-1].value != "(" {
					output = append(output, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
				if len(stack) == 0 {
					return nil, errors.New("mismatched parentheses")
				}
				stack = stack[:len(stack)-1]
			} else {
				for len(stack) > 0 && precedence[stack[len(stack)-1].value] >= precedence[token.value] {
					output = append(output, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, token)
			}
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1].value == "(" || stack[len(stack)-1].value == ")" {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

// evaluatePostfix evaluates a postfix expression
func evaluatePostfix(tokens []Token) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		if token.typ == "number" {
			number, err := strconv.ParseFloat(token.value, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, number)
		} else if token.typ == "operator" {
			if len(stack) < 2 {
				return 0, errors.New("insufficient values")
			}
			right, left := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token.value {
			case "+":
				result = left + right
			case "-":
				result = left - right
			case "*":
				result = left * right
			case "/":
				if right == 0 {
					return 0, errors.New("division by zero")
				}
				result = left / right
			case "^":
				result = math.Pow(left, right)
			default:
				return 0, errors.New("unknown operator")
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}
