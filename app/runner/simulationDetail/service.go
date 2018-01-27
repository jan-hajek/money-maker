package simulationDetail

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/dateInput"
	"github.com/jelito/money-maker/app/log"
	"github.com/jelito/money-maker/app/registry"
)

type Service struct {
	Log             log.Log
	Writer          *app.Writer
	Strategies      map[string]map[string]map[string]interface{}
	Registry        *registry.Registry
	DateInputLoader dateInput.Loader
}

func (s *Service) Run() {

	strategies := s.loadStrategies()

	s.Log.Info("strategies: ", len(strategies))

	dateInputs, err := s.DateInputLoader.Load()
	if err != nil {
		s.Log.Fatal(err)
	}

	err = s.Writer.Open()
	if err != nil {
		s.Log.Fatal(err)
	}

	for _, strategy := range strategies {
		history := s.runStrategy(strategy, dateInputs)
		summary := app.CreateSummary(history)

		err = s.Writer.WriteHistory(history.GetAll())
		if err != nil {
			s.Log.Fatal(err)
		}

		err = s.Writer.WriteSummaryHeader(summary)
		if err != nil {
			s.Log.Fatal(err)
		}
		err = s.Writer.WriteSummaryRow(summary)
		if err != nil {
			s.Log.Fatal(err)
		}
	}

	err = s.Writer.Close()
	if err != nil {
		s.Log.Fatal(err)
	}

}

func (s *Service) loadStrategies() []app.Strategy {
	var strategies []app.Strategy

	for strategyFactoryName, strategyFactoryConfig := range s.Strategies {
		strategyFactory := s.Registry.GetByName(strategyFactoryName).(app.StrategyFactory)

		strategies = append(
			strategies,
			strategyFactory.Create(
				strategyFactory.GetDefaultConfig(strategyFactoryConfig),
			),
		)
	}
	return strategies
}

func (s *Service) runStrategy(strategy app.Strategy, dateInputs []app.DateInput) *app.History {
	var lastPosition *app.Position
	indicators := strategy.GetIndicators()

	history := &app.History{
		Strategy:   strategy,
		Indicators: indicators,
	}

	iteration := 0
	for _, dateInput := range dateInputs {

		indicatorResults := map[string]app.IndicatorResult{}

		iteration++

		for _, c := range indicators {
			input := app.IndicatorInput{
				Date:       dateInput.Date,
				OpenPrice:  dateInput.OpenPrice,
				HighPrice:  dateInput.HighPrice,
				LowPrice:   dateInput.LowPrice,
				ClosePrice: dateInput.ClosePrice,
				Iteration:  iteration,
			}

			indicatorResults[c.GetName()] = c.Calculate(input, history)
		}

		strategyResult := strategy.Resolve(app.StrategyInput{
			DateInput:        dateInput,
			History:          history,
			Position:         lastPosition,
			IndicatorResults: indicatorResults,
		})

		lastPosition = createPosition(strategyResult, dateInput, lastPosition)

		historyItem := &app.HistoryItem{
			DateInput:        dateInput,
			IndicatorResults: indicatorResults,
			StrategyResult:   strategyResult,
			Position:         lastPosition,
		}

		history.AddItem(historyItem)

		if strategyResult.Action == app.CLOSE {
			lastPosition = nil
		}
	}

	return history
}
