package trade

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/log"
	"github.com/jelito/money-maker/app/repository/position"
)

type Service struct {
	positionRepository *position.Service

	Trade        *entity.Trade
	strategy     app.Strategy
	indicators   []app.Indicator
	iteration    int
	history      *app.History
	lastPosition *app.Position
	log          log.Log
}

func (s *Service) Run(dateInput app.DateInput) (*app.History, error) {
	s.iteration++

	var err error
	s.lastPosition, err = s.getLastPosition()
	if err != nil {
		return nil, err
	}

	s.history.SetLastPosition(s.lastPosition)

	indicatorResults := s.getIndicatorResults(dateInput)

	strategyResult := s.strategy.Resolve(app.StrategyInput{
		DateInput:        dateInput,
		History:          s.history,
		Position:         s.lastPosition,
		IndicatorResults: indicatorResults,
	})

	historyItem := &app.HistoryItem{
		DateInput:        dateInput,
		IndicatorResults: indicatorResults,
		StrategyResult:   strategyResult,
	}
	s.history.AddItem(historyItem)

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

func (s *Service) getLastPosition() (*app.Position, error) {

	positionEntity, err := s.positionRepository.LastOpenByTrade(s.Trade.Id)
	if err != nil {
		return nil, err
	}
	if positionEntity == nil {
		return nil, nil
	}

	return &app.Position{
		Id:     positionEntity.Id,
		Type:   app.PositionType(positionEntity.Type),
		Amount: float.New(positionEntity.Amount),
		Sl:     float.New(positionEntity.Sl),
		Costs:  float.New(positionEntity.Costs),
	}, nil

}
