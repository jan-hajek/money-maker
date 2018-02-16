package adxEma

import (
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/printValue"
)

type Result struct {
	Adx     float.Float
	DIPlus  float.Float
	DIMinus float.Float
}

func (s Result) Print() []printValue.PrintValue {
	return []printValue.PrintValue{
		{Label: "adx", Value: s.Adx},
		{Label: "diPlus", Value: s.DIPlus},
		{Label: "diMinus", Value: s.DIMinus},
	}
}
