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

func (s *Position) clone() *Position {
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

var lastPositionId = 0

func createPosition(resolverResult ResolverResult, dateInput DateInput, lastPosition *Position) *Position {

	switch resolverResult.Action {
	case OPEN:
		if lastPosition != nil {
			panic("position is already open")
		}
		lastPositionId++
		newPosition := &Position{
			Id:        lastPositionId,
			Type:      resolverResult.PositionType,
			StartDate: dateInput.Date,
			OpenPrice: dateInput.ClosePrice,
			Amount:    resolverResult.Amount,
			Sl:        resolverResult.Sl,
			Costs:     resolverResult.Costs,
		}

		newPosition.Profit = calculateProfit(newPosition)
		newPosition.PossibleProfit = calculatePossibleProfit(newPosition, dateInput.ClosePrice)

		return newPosition
	case CLOSE:
		if lastPosition == nil {
			panic("you can't close, position is not open")
		}
		newPosition := lastPosition.clone()

		newPosition.ClosePrice = dateInput.ClosePrice
		newPosition.Sl = resolverResult.Sl
		newPosition.Costs = resolverResult.Costs
		newPosition.Profit = calculateProfit(newPosition)
		newPosition.PossibleProfit = newPosition.Profit

		return newPosition
	case CHANGE:
		if lastPosition == nil {
			panic("you can't change, position is not open")
		}
		newPosition := lastPosition.clone()

		newPosition.Sl = resolverResult.Sl
		newPosition.Costs = resolverResult.Costs
		newPosition.Profit = calculateProfit(newPosition)
		newPosition.PossibleProfit = calculatePossibleProfit(newPosition, dateInput.ClosePrice)

		return newPosition
	case SKIP:
		if lastPosition != nil {
			newPosition := lastPosition.clone()

			newPosition.Profit = calculateProfit(newPosition)
			newPosition.PossibleProfit = calculatePossibleProfit(newPosition, dateInput.ClosePrice)

			return newPosition
		}
	}

	return lastPosition
}

func calculateProfit(position *Position) float.Float {
	profit := position.Amount.Val() * (position.ClosePrice.Val() - position.OpenPrice.Val())
	if position.Type == SHORT {
		profit *= -1
	}

	return float.New(profit - position.Costs.Val())
}

func calculatePossibleProfit(position *Position, actualPrice float.Float) float.Float {
	profit := position.Amount.Val() * (actualPrice.Val() - position.OpenPrice.Val())
	if position.Type == SHORT {
		profit *= -1
	}

	return float.New(profit - position.Costs.Val())
}

type PositionType string

const (
	LONG  PositionType = "long"
	SHORT              = "short"
)
