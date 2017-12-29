package samson

type Config struct {
	SarMinimalAf float64
	SarMaximalAf float64
	AdxPeriod    int
	SmoothType   SmoothType
}

type SmoothType string

const (
	AVG SmoothType = "avg"
	EMA            = "ema"
)
