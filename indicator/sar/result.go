package sar

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
)

type Result struct {
	Sar     float.Float
	Af      float.Float
	Ep      float.Float
	UpTrend bool
}

func (s Result) Print() []app.PrintValue {
	return []app.PrintValue{
		{Label: "sar", Value: s.Sar},
	}
}
