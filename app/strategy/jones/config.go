package jones

import "github.com/jelito/money-maker/src/math/float"

type Config struct {
	AdxPeriod        int
	SmoothType       SmoothType
	SmoothAlpha      float.Float
	OpenLowerAdxVal  int
	OpenHigherAdxVal int
	CloseAdxVal      int
}

type SmoothType string

const (
	AVG SmoothType = "avg"
	EMA            = "ema"
)
