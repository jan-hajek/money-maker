package position

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/satori/go.uuid"
)

type Maker struct {
}

func (s *Maker) Create(
	strategyResult app.StrategyResult,
	dateInput app.DateInput,
	lastPosition *app.Position,
) *app.Position {

	switch strategyResult.Action {
	case app.OPEN:
		if lastPosition != nil {
			panic("position is already open")
		}
		newPosition := &app.Position{
			Id:        uuid.Must(uuid.NewV4()).String(),
			Type:      strategyResult.PositionType,
			StartDate: dateInput.Date,
			OpenPrice: dateInput.ClosePrice,
			Amount:    strategyResult.Amount,
			Sl:        strategyResult.Sl,
			Costs:     strategyResult.Costs,
		}

		newPosition.Profit = s.calculateProfit(newPosition)
		newPosition.PossibleProfit = s.calculatePossibleProfit(newPosition, dateInput.ClosePrice)

		return newPosition
	case app.CLOSE:
		if lastPosition == nil {
			panic("you can't close, position is not open")
		}
		newPosition := lastPosition.Clone()

		newPosition.ClosePrice = dateInput.ClosePrice
		newPosition.Sl = strategyResult.Sl
		newPosition.Costs = strategyResult.Costs
		newPosition.Profit = s.calculateProfit(newPosition)
		newPosition.PossibleProfit = newPosition.Profit

		return newPosition
	case app.CHANGE:
		if lastPosition == nil {
			panic("you can't change, position is not open")
		}
		newPosition := lastPosition.Clone()

		newPosition.Sl = strategyResult.Sl
		newPosition.Costs = strategyResult.Costs
		newPosition.Profit = s.calculateProfit(newPosition)
		newPosition.PossibleProfit = s.calculatePossibleProfit(newPosition, dateInput.ClosePrice)

		return newPosition
	case app.SKIP:
		if lastPosition != nil {
			newPosition := lastPosition.Clone()

			newPosition.Costs = strategyResult.Costs
			newPosition.Profit = s.calculateProfit(newPosition)
			newPosition.PossibleProfit = s.calculatePossibleProfit(newPosition, dateInput.ClosePrice)

			return newPosition
		}
	}

	return lastPosition
}

func (s *Maker) calculateProfit(position *app.Position) float.Float {
	profit := position.Amount.Val() * (position.ClosePrice.Val() - position.OpenPrice.Val())
	if position.Type == app.SHORT {
		profit *= -1
	}

	return float.New(profit)
}

func (s *Maker) calculatePossibleProfit(position *app.Position, actualPrice float.Float) float.Float {
	profit := position.Amount.Val() * (actualPrice.Val() - position.OpenPrice.Val())
	if position.Type == app.SHORT {
		profit *= -1
	}

	return float.New(profit)
}
