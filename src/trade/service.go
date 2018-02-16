package trade

import (
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/entity"
	"github.com/jelito/money-maker/src/history"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/log"
	"github.com/jelito/money-maker/src/position"
	"github.com/jelito/money-maker/src/strategy"
)

type Service struct {
	Trade         *entity.Trade
	strategy      interfaces.Strategy
	indicators    []interfaces.Indicator
	history       interfaces.History
	positionMaker *position.PositionMaker
	log           log.Log
	idGenerator   position.IdGenerator

	iteration    int
	lastPosition *interfaces.Position
}

func (s *Service) Run(dateInput dateInput.DateInput) (interfaces.History, error) {
	s.iteration++

	indicatorResults := s.getIndicatorResults(dateInput)

	strategyResult := s.strategy.Resolve(&strategy.Input{
		DateInput:        dateInput,
		History:          s.history,
		Position:         s.lastPosition,
		IndicatorResults: indicatorResults,
	})

	s.lastPosition = s.positionMaker.Create(strategyResult, dateInput, s.lastPosition, s.idGenerator)

	historyItem := &history.HistoryItem{
		DateInput:        dateInput,
		IndicatorResults: indicatorResults,
		StrategyResult:   strategyResult,
		Position:         s.lastPosition,
	}
	s.history.AddItem(historyItem)

	if strategyResult.GetAction() == interfaces.CLOSE {
		s.lastPosition = nil
	}

	return s.history, nil
}

func (s *Service) getIndicatorResults(dateInput dateInput.DateInput) map[string]interfaces.IndicatorResult {
	indicatorResults := map[string]interfaces.IndicatorResult{}

	// TODO - jhajek rutiny
	for _, c := range s.indicators {
		input := interfaces.IndicatorInput{
			Date:       dateInput.Date,
			OpenPrice:  dateInput.OpenPrice,
			HighPrice:  dateInput.HighPrice,
			LowPrice:   dateInput.LowPrice,
			ClosePrice: dateInput.ClosePrice,
			Iteration:  s.iteration,
		}

		indicatorResults[c.GetName()] = c.Calculate(input, s.history)
	}

	return indicatorResults
}
