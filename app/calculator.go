package app

import (
	"github.com/jelito/money-maker/app/float"
	"time"
)

type Indicator interface {
	Calculate(input IndicatorInput, history *History) IndicatorResult
	GetName() string
}

type IndicatorInput struct {
	Date       time.Time
	OpenPrice  float.Float
	ClosePrice float.Float
	HighPrice  float.Float
	LowPrice   float.Float
	Iteration  int
}

type IndicatorResult interface {
	Print() []PrintValue
}
