package position

import (
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math/float"
)

type PositionMaker struct {
}

func (s *PositionMaker) Create(
	strategyResult interfaces.StrategyResult,
	dateInput dateInput.DateInput,
	lastPosition *interfaces.Position,
	idGenerator IdGenerator,
) *interfaces.Position {
	switch strategyResult.GetAction() {
	case interfaces.OPEN:
		if lastPosition != nil {
			panic("position is already open")
		}
		newPosition := &interfaces.Position{
			Id:        idGenerator.Generate(),
			Type:      strategyResult.GetPositionType(),
			StartDate: dateInput.Date,
			OpenPrice: dateInput.ClosePrice,
			Amount:    strategyResult.GetAmount(),
			Sl:        strategyResult.GetSl(),
			Costs:     strategyResult.GetCosts(),
		}

		newPosition.PossibleProfitPercent = s.calculatePossibleProfitPercent(newPosition, dateInput.ClosePrice)
		newPosition.PossibleProfit = s.calculatePossibleProfit(newPosition)

		return newPosition
	case interfaces.CLOSE:
		if lastPosition == nil {
			panic("you can't close, position is not open")
		}
		newPosition := lastPosition.Clone()

		newPosition.ClosePrice = dateInput.ClosePrice
		newPosition.Sl = strategyResult.GetSl()
		newPosition.Costs = strategyResult.GetCosts()

		newPosition.Profit = s.calculateProfit(newPosition)
		newPosition.PossibleProfitPercent = s.calculatePossibleProfitPercent(newPosition, dateInput.ClosePrice)
		newPosition.PossibleProfit = s.calculatePossibleProfit(newPosition)

		return newPosition
	case interfaces.CHANGE:
		if lastPosition == nil {
			panic("you can't change, position is not open")
		}
		newPosition := lastPosition.Clone()

		newPosition.Sl = strategyResult.GetSl()
		newPosition.Costs = strategyResult.GetCosts()

		newPosition.PossibleProfitPercent = s.calculatePossibleProfitPercent(newPosition, dateInput.ClosePrice)
		newPosition.PossibleProfit = s.calculatePossibleProfit(newPosition)

		return newPosition
	case interfaces.SKIP:
		if lastPosition != nil {
			newPosition := lastPosition.Clone()

			newPosition.Costs = strategyResult.GetCosts()

			newPosition.PossibleProfitPercent = s.calculatePossibleProfitPercent(newPosition, dateInput.ClosePrice)
			newPosition.PossibleProfit = s.calculatePossibleProfit(newPosition)

			return newPosition
		}
	}

	return lastPosition
}

func (s *PositionMaker) calculateProfit(p *interfaces.Position) float.Float {
	profit := p.Amount.Val() * (p.ClosePrice.Val() - p.OpenPrice.Val())
	if p.Type == interfaces.SHORT {
		profit *= -1
	}

	return float.New(profit)
}

func (s *PositionMaker) calculatePossibleProfitPercent(p *interfaces.Position, actualPrice float.Float) float.Float {
	percent := actualPrice.Sub(p.OpenPrice).Div(p.OpenPrice)
	if p.Type == interfaces.SHORT {
		percent = percent.MultiFloat(-1)
	}

	return percent
}

func (s *PositionMaker) calculatePossibleProfit(p *interfaces.Position) float.Float {
	return p.Amount.Multi(p.OpenPrice).Multi(p.PossibleProfitPercent)
}
