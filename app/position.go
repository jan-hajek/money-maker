package app

import (
	"github.com/jelito/money-maker/app/float"
	"time"
)

type Position struct {
	Id             int
	Type           PositionType
	StartDate      time.Time
	CloseDate      time.Time
	OpenPrice      float.Float
	ClosePrice     float.Float
	Amount         float.Float
	Sl             float.Float
	Costs          float.Float
	Profit         float.Float
	PossibleProfit float.Float
}

func (s *Position) Clone() *Position {
	return &Position{
		Id:             s.Id,
		Type:           s.Type,
		StartDate:      s.StartDate,
		CloseDate:      s.CloseDate,
		OpenPrice:      s.OpenPrice,
		ClosePrice:     s.ClosePrice,
		Amount:         s.Amount,
		Sl:             s.Sl,
		Costs:          s.Costs,
		Profit:         s.Profit,
		PossibleProfit: s.PossibleProfit,
	}
}

type PositionType string

const (
	LONG  PositionType = "long"
	SHORT              = "short"
)
