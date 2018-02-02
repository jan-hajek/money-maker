package simulationBatch

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/dateInput"
	"github.com/jelito/money-maker/app/log"
	"github.com/jelito/money-maker/app/registry"
	"github.com/jelito/money-maker/app/writer"
)

type Service struct {
	Log             log.Log
	Writer          *writer.Writer
	Strategies      map[string]map[string]map[string]interface{}
	Registry        *registry.Registry
	DateInputLoader dateInput.Loader
}

func (s *Service) Run() {

	strategies := s.loadStrategies()

	uiprogress.Start()
	bar := uiprogress.AddBar(len(strategies)).AppendCompleted().AppendElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("(%d/%d)", b.Current(), b.Total)
	})

	dateInputs, err := s.DateInputLoader.Load()
	if err != nil {
		s.Log.Fatal(err)
	}

	err = s.Writer.Open()
	if err != nil {
		s.Log.Fatal(err)
	}

	ch := make(chan *app.Summary)
	x := make(chan int, 8)
	wait := make(chan bool)

	go func() {
		for i := 0; i < len(strategies); i++ {
			summary := <-ch
			<-x

			if i == 0 {

				err = s.Writer.WriteSummaryHeader(summary)
				if err != nil {
					s.Log.Fatal(err)
				}
			}

			err = s.Writer.WriteSummaryRow(summary)
			if err != nil {
				s.Log.Fatal(err)

			}

			bar.Incr()
		}

		uiprogress.Stop()

		wait <- true
	}()

	for index, strategy := range strategies {
		go func(index int, strategy app.Strategy) {
			history := s.runStrategy(strategy, dateInputs)

			ch <- app.CreateSummary(history)
		}(index, strategy)

		x <- index

	}

	<-wait

	err = s.Writer.Close()
	if err != nil {
		s.Log.Fatal(err)
	}

	s.Log.Info("done")

}

func (s *Service) loadStrategies() []app.Strategy {
	var strategies []app.Strategy

	for strategyFactoryName, strategyFactoryConfig := range s.Strategies {
		strategyFactory := s.Registry.GetByName(strategyFactoryName).(app.StrategyFactory)

		for _, config := range strategyFactory.GetBatchConfigs(strategyFactoryConfig) {
			strategies = append(strategies, strategyFactory.Create(config))
		}
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
