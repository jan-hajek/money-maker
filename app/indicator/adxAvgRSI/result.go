package adxAvgRSI

import (
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/printValue"
)

type Result struct {
	Adx           float.Float
	TrueRange     float.Float
	EmaTrueRange  float.Float
	DmPlus        float.Float
	EmaDmPlus     float.Float
	DmMinus       float.Float
	EmaDmMinus    float.Float
	DIAbs         float.Float
	SmmaDIAbs     float.Float
	DIPlus        float.Float
	DIMinus       float.Float
	AdxDiff       float.Float
	AdxGain       float.Float
	AdxLoss       float.Float
	SmoothAdxGain float.Float
	SmoothAdxLoss float.Float
	GainToLoss    float.Float
	DX            float.Float
	RSI           float.Float
}

func (s Result) Print() []printValue.PrintValue {
	return []printValue.PrintValue{
		{Label: "adx", Value: s.Adx},
		{Label: "diPlus", Value: s.DIPlus},
		{Label: "diMinus", Value: s.DIMinus},
		{Label: "DX", Value: s.DX},
		{Label: "RSI", Value: s.RSI},
	}
}
