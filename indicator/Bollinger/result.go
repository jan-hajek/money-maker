package Bollinger

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
)

type Result struct {
	TypicalPrice float.Float
	BL           float.Float
	BMA          float.Float
	BU           float.Float
}

func (s Result) Print() []app.PrintValue {
	return []app.PrintValue{
		{Label: "BL", Value: s.BL},
		{Label: "BMA", Value: s.BMA},
		{Label: "BU", Value: s.BU},
	}
}
