package app

import (
	"github.com/jelito/money-maker/app/float"
	"time"
)

type Calculator interface {
	Calculate(input CalculatorInput, history *History) CalculatorResult
	GetName() string
}

type CalculatorInput struct {
	Date       time.Time
	OpenPrice  float.Float
	ClosePrice float.Float
	HighPrice  float.Float
	LowPrice   float.Float
	Iteration  int
}

type CalculatorResult interface {
	Print() []PrintValue
}
