package runner

import (
	"github.com/jelito/money-maker/app"
	"log"
)

func (s App) Run() {

	strategies := s.loadStrategies(
		func(
			strategyFactory app.StrategyFactory,
			config map[string]map[string]interface{},
		) []app.StrategyFactoryConfig {

			return []app.StrategyFactoryConfig{
				strategyFactory.GetDefaultConfig(config),
			}
		},
	)

	log.Println("strategies: ", len(strategies))

	dateInputs, err := getDateInputs(s.Config.InputFile, s.Config.ParseFormat)
	if err != nil {
		log.Fatal(err)
	}

	writer := s.createWriter()
	err = writer.Open()
	if err != nil {
		log.Fatal(err)
	}

	for _, strategy := range strategies {
		history := s.runStrategy(strategy, dateInputs)
		summary := app.CreateSummary(history)

		err = writer.WriteHistory(history)
		if err != nil {
			log.Fatal(err)
		}

		err = writer.WriteSummaryHeader(summary)
		if err != nil {
			log.Fatal(err)
		}
		err = writer.WriteSummaryRow(summary)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}
}
