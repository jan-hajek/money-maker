package simulationDetail

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"strconv"
)

var lastPositionId = 0

func createPosition(strategyResult app.StrategyResult, dateInput app.DateInput, lastPosition *app.Position) *app.Position {

	switch strategyResult.Action {
	case app.OPEN:
		if lastPosition != nil {
			panic("position is already open")
		}
		lastPositionId++
		newPosition := &app.Position{
			Id:        strconv.Itoa(lastPositionId),
			Type:      strategyResult.PositionType,
			StartDate: dateInput.Date,
			OpenPrice: dateInput.ClosePrice,
			Amount:    strategyResult.Amount,
			Sl:        strategyResult.Sl,
			Costs:     strategyResult.Costs,
		}

		newPosition.Profit = calculateProfit(newPosition)
		newPosition.PossibleProfit = calculatePossibleProfit(newPosition, dateInput.ClosePrice)

		return newPosition
	case app.CLOSE:
		if lastPosition == nil {
			panic("you can't close, position is not open")
		}
		newPosition := lastPosition.Clone()

		newPosition.ClosePrice = dateInput.ClosePrice
		newPosition.Sl = strategyResult.Sl
		newPosition.Costs = strategyResult.Costs
		newPosition.Profit = calculateProfit(newPosition)
		newPosition.PossibleProfit = newPosition.Profit

		return newPosition
	case app.CHANGE:
		if lastPosition == nil {
			panic("you can't change, position is not open")
		}
		newPosition := lastPosition.Clone()

		newPosition.Sl = strategyResult.Sl
		newPosition.Costs = strategyResult.Costs
		newPosition.Profit = calculateProfit(newPosition)
		newPosition.PossibleProfit = calculatePossibleProfit(newPosition, dateInput.ClosePrice)

		return newPosition
	case app.SKIP:
		if lastPosition != nil {
			newPosition := lastPosition.Clone()

			newPosition.Profit = calculateProfit(newPosition)
			newPosition.PossibleProfit = calculatePossibleProfit(newPosition, dateInput.ClosePrice)

			return newPosition
		}
	}

	return lastPosition
}

func calculateProfit(position *app.Position) float.Float {
	profit := position.Amount.Val() * (position.ClosePrice.Val() - position.OpenPrice.Val())
	if position.Type == app.SHORT {
		profit *= -1
	}

	return float.New(profit)
}

func calculatePossibleProfit(position *app.Position, actualPrice float.Float) float.Float {
	profit := position.Amount.Val() * (actualPrice.Val() - position.OpenPrice.Val())
	if position.Type == app.SHORT {
		profit *= -1
	}

	return float.New(profit)
}
