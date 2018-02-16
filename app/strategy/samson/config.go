package samson

import "github.com/jelito/money-maker/src/math/float"

type Config struct {
	SarMinimalAf float.Float
	SarMaximalAf float.Float
	AdxPeriod    int
	SmoothType   SmoothType
	SmoothAlpha  float.Float
}

type SmoothType string

const (
	AVG SmoothType = "avg"
	EMA            = "ema"
)
