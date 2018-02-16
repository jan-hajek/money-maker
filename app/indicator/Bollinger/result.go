package Bollinger

import (
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/printValue"
)

type Result struct {
	TypicalPrice float.Float
	BL           float.Float
	BMA          float.Float
	BU           float.Float
}

func (s Result) Print() []printValue.PrintValue {
	return []printValue.PrintValue{
		{Label: "BL", Value: s.BL},
		{Label: "BMA", Value: s.BMA},
		{Label: "BU", Value: s.BU},
	}
}
