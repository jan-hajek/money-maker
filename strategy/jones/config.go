package jones

import "github.com/jelito/money-maker/app/float"

type Config struct {
	AdxPeriod   int
	SmoothType  SmoothType
	SmoothAlpha float.Float
}

type SmoothType string

const (
	AVG SmoothType = "avg"
	EMA            = "ema"
)
