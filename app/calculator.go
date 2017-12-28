package app

import (
	"time"
)

type Calculator interface {
	Calculate(input CalculatorInput, history *History) CalculatorResult
	GetName() string
}

type CalculatorInput struct {
	Date       time.Time
	OpenPrice  float64
	ClosePrice float64
	HighPrice  float64
	LowPrice   float64
	Iteration  int
}

type CalculatorResult interface {
	Print() []PrintValue
}
