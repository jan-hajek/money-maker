package app

import "time"

type DateInput struct {
	Date       time.Time
	OpenPrice  float64
	ClosePrice float64
	HighPrice  float64
	LowPrice   float64
}
