package runner

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/registry"
)

type App struct {
	Config   Config
	Registry *registry.Registry
}

// TODO - jhajek yzanoreni
type Config struct {
	InputFile   string                                       `yaml:"input"`
	ParseFormat string                                       `yaml:"parseFormat"`
	Outputs     map[string]interface{}                       `yaml:"outputs"`
	Strategies  map[string]map[string]map[string]interface{} `yaml:"strategies"`
}

func (s *App) loadStrategies(
	loadConfigFunc func(app.StrategyFactory, map[string]map[string]interface{}) []app.StrategyFactoryConfig,
) []app.Strategy {
	var strategies []app.Strategy

	for strategyFactoryName, strategyFactoryConfig := range s.Config.Strategies {
		strategyFactory := s.Registry.GetByName(strategyFactoryName).(app.StrategyFactory)

		for _, config := range loadConfigFunc(strategyFactory, strategyFactoryConfig) {
			strategies = append(strategies, strategyFactory.Create(config))
		}
	}
	return strategies
}

func (s *App) runStrategy(strategy app.Strategy, dateInputs []app.DateInput) *app.History {
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

func (s *App) createWriter() app.Writer {
	var outputs []app.WriterOutput

	if value, ok := s.Config.Outputs["stdout"]; ok == true && value == true {
		outputs = append(outputs, &app.StdOutWriterOutput{
			DateFormat: s.Config.ParseFormat,
		})
	}

	if value, ok := s.Config.Outputs["csv"]; ok == true && value != "" {
		outputs = append(outputs, &app.CsvWriterOutput{
			File:       value.(string),
			DateFormat: s.Config.ParseFormat,
		})
	}

	return app.Writer{
		Outputs: outputs,
	}
}
