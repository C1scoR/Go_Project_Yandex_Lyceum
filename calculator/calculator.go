package calculator

import (
	"errors"
	"strconv"
	"strings"
)

// Calc вычисляет результат арифметического выражения
func Calc(expression string) (float64, error) {
	// Удаляем пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")

	var stack []float64  // стек для хранения чисел
	var operators []rune // стек для хранения операторов

	priority := map[rune]int{'+': 1, '-': 1, '*': 2, '/': 2}

	// Функция для выполнения операций
	applyOperator := func() error {
		if len(stack) < 2 || len(operators) == 0 {
			return errors.New("некорректное выражение")
		}
		b := stack[len(stack)-1]
		a := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		var result float64
		switch op {
		case '+':
			result = a + b
		case '-':
			result = a - b
		case '*':
			result = a * b
		case '/':
			if b == 0 {
				return errors.New("деление на ноль")
			}
			result = a / b
		}
		stack = append(stack, result)
		return nil
	}

	// Разбор выражения
	for i := 0; i < len(expression); i++ {
		switch {
		case expression[i] >= '0' && expression[i] <= '9':
			// Если это цифра, считываем число
			start := i
			for i+1 < len(expression) && (expression[i+1] >= '0' && expression[i+1] <= '9' || expression[i+1] == '.') {
				i++
			}
			num, err := strconv.ParseFloat(expression[start:i+1], 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		case expression[i] == '(':
			operators = append(operators, '(')
		case expression[i] == ')':
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				if err := applyOperator(); err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 {
				return 0, errors.New("некорректное выражение: недостаточно открывающих скобок")
			}
			operators = operators[:len(operators)-1] // удаляем '('
		default:
			for len(operators) > 0 && priority[rune(expression[i])] <= priority[operators[len(operators)-1]] {
				if err := applyOperator(); err != nil {
					return 0, err
				}
			}
			operators = append(operators, rune(expression[i]))
		}
	}

	// Применяем оставшиеся операторы
	for len(operators) > 0 {
		if err := applyOperator(); err != nil {
			return 0, err
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("некорректное выражение")
	}

	return stack[0], nil
}
