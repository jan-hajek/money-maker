package trade

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/log"
	"github.com/jelito/money-maker/app/position"
)

type Service struct {
	Trade         *entity.Trade
	strategy      app.Strategy
	indicators    []app.Indicator
	history       *app.History
	positionMaker *position.Maker
	log           log.Log

	iteration    int
	lastPosition *app.Position
}

func (s *Service) Run(dateInput app.DateInput) (*app.History, error) {
	s.iteration++

	indicatorResults := s.getIndicatorResults(dateInput)

	strategyResult := s.strategy.Resolve(app.StrategyInput{
		DateInput:        dateInput,
		History:          s.history,
		Position:         s.lastPosition,
		IndicatorResults: indicatorResults,
	})

	s.lastPosition = s.positionMaker.Create(strategyResult, dateInput, s.lastPosition)

	historyItem := &app.HistoryItem{
		DateInput:        dateInput,
		IndicatorResults: indicatorResults,
		StrategyResult:   strategyResult,
		Position:         s.lastPosition,
	}
	s.history.AddItem(historyItem)

	if strategyResult.Action == app.CLOSE {
		s.lastPosition = nil
	}

	return s.history, nil
}

func (s *Service) getIndicatorResults(dateInput app.DateInput) map[string]app.IndicatorResult {
	indicatorResults := map[string]app.IndicatorResult{}

	// TODO - jhajek rutiny
	for _, c := range s.indicators {
		input := app.IndicatorInput{
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
