package adxEmaRSI

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
)

type Result struct {
	Adx          float.Float
	TrueRange    float.Float
	EmaTrueRange float.Float
	DmPlus       float.Float
	EmaDmPlus    float.Float
	DmMinus      float.Float
	EmaDmMinus   float.Float
	DIAbs        float.Float
	SmmaDIAbs    float.Float
	DIPlus       float.Float
	DIMinus      float.Float
	DIDiff       float.Float
	DIDiff3      float.Float
	DIDiffLB     float.Float
	DIDiffMA     float.Float
	DIDiffUB     float.Float
}

func (s Result) Print() []app.PrintValue {
	return []app.PrintValue{
		{Label: "adx", Value: s.Adx},
		{Label: "diPlus", Value: s.DIPlus},
		{Label: "diMinus", Value: s.DIMinus},
		{Label: "DIDiff3", Value: s.DIDiff3},
		//{Label: "DIDiffLB", Value: s.DIDiffLB},
		//{Label: "DIDiffMA", Value: s.DIDiffMA},
		//{Label: "DIDiffUB", Value: s.DIDiffUB},
	}
}
