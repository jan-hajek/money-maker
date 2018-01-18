package entity

import (
	"time"
)

type Price struct {
	Id         string
	TitleId    string
	Date       time.Time
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
}
