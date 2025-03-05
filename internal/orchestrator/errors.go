package orchestrator

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDivisionByZero    = errors.New("division by zero")
	ErrDotEndOfOperand   = errors.New("dot at the end of the operand")
)
