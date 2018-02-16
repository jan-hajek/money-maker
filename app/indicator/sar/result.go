package sar

import (
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/printValue"
)

type Result struct {
	Sar     float.Float
	Af      float.Float
	Ep      float.Float
	UpTrend bool
}

func (s Result) Print() []printValue.PrintValue {
	return []printValue.PrintValue{
		{Label: "sar", Value: s.Sar},
	}
}
