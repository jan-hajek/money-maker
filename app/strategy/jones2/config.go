package jones2

import (
	"github.com/jelito/money-maker/src/math/float"
)

type Config struct {
	AdxPeriod     int
	SmoothType    SmoothType
	SmoothAlpha   float.Float
	OpenLowerAdx  int
	OpenHigherAdx int
	CloseAdx      int
	DIOpenLevel   int
	DICloseLevel  int
	DISDCount     float.Float
	PeriodDIMA    int
	Spread        float.Float
	Swap          float.Float
	StopProfit    float.Float
}

type SmoothType string

const (
	AVG SmoothType = "avg"
	EMA            = "ema"
)
