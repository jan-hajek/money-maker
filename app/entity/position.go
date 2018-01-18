package entity

import (
	"database/sql"
)

type Position struct {
	Id           string
	TradeId      string
	Type         string
	OpenPriceId  string
	ClosePriceId sql.NullString
	Amount       float64
	Sl           float64
	Costs        float64
	Profit       float64
}
