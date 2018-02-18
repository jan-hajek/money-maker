package simulationBatch

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/history"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/log"
	"github.com/jelito/money-maker/src/position"
	"github.com/jelito/money-maker/src/registry"
	"github.com/jelito/money-maker/src/strategy"
	"github.com/jelito/money-maker/src/summary"
	"github.com/jelito/money-maker/src/writer"
)

type Service struct {
	Log             log.Log
	Writer          *writer.Writer
	Strategies      map[string]map[string]map[string]interface{}
	Registry        *registry.Registry
	DateInputLoader dateInput.Loader
	PositionMaker   *position.PositionMaker
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

	ch := make(chan *summary.Summary)
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
		go func(index int, strategy interfaces.Strategy) {
			history := s.runStrategy(strategy, dateInputs)

			ch <- summary.CreateSummary(history)
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

func (s *Service) loadStrategies() []interfaces.Strategy {
	var strategies []interfaces.Strategy

	for strategyFactoryName, strategyFactoryConfig := range s.Strategies {
		strategyFactory := s.Registry.GetByName(strategyFactoryName).(strategy.Factory)

		for _, config := range strategyFactory.GetBatchConfigs(strategyFactoryConfig) {
			strategies = append(strategies, strategyFactory.Create(config))
		}
	}
	return strategies
}

func (s *Service) runStrategy(strategyItem interfaces.Strategy, dateInputs []dateInput.DateInput) interfaces.History {
	var lastPosition *interfaces.Position
	indicators := strategyItem.GetIndicators()

	idGenerator := &position.UuidGenerator{}

	historyItem := &history.History{
		Strategy:   strategyItem,
		Indicators: indicators,
	}

	iteration := 0
	for _, dateInput := range dateInputs {

		indicatorResults := map[string]interfaces.IndicatorResult{}

		iteration++

		for _, c := range indicators {
			input := interfaces.IndicatorInput{
				Date:       dateInput.Date,
				OpenPrice:  dateInput.OpenPrice,
				HighPrice:  dateInput.HighPrice,
				LowPrice:   dateInput.LowPrice,
				ClosePrice: dateInput.ClosePrice,
				Iteration:  iteration,
			}

			indicatorResults[c.GetName()] = c.Calculate(input, historyItem)
		}

		strategyResult := strategyItem.Resolve(&strategy.Input{
			DateInput:        dateInput,
			History:          historyItem,
			Position:         lastPosition,
			IndicatorResults: indicatorResults,
		})

		lastPosition = s.PositionMaker.Create(strategyResult, dateInput, lastPosition, idGenerator)

		historyItem.AddItem(&history.HistoryItem{
			DateInput:        dateInput,
			IndicatorResults: indicatorResults,
			StrategyResult:   strategyResult,
			Position:         lastPosition,
		})

		if strategyResult.GetAction() == interfaces.CLOSE {
			lastPosition = nil
		}
	}

	return historyItem
}
