package app

import (
	"github.com/jelito/money-maker/app/float"
	"time"
)

type DateInput struct {
	Date       time.Time
	OpenPrice  float.Float
	ClosePrice float.Float
	HighPrice  float.Float
	LowPrice   float.Float
}
