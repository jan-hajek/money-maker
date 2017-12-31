package runner

import (
	"github.com/jelito/money-maker/app"
	"gopkg.in/cheggaaa/pb.v1"
	"log"
)

func (s App) Batch() {

	strategies := s.loadStrategies(
		func(
			strategyFactory app.StrategyFactory,
			config map[string]map[string]interface{},
		) []app.StrategyFactoryConfig {

			return strategyFactory.GetBatchConfigs(config)

		},
	)

	bar := pb.StartNew(len(strategies))

	dateInputs, err := getDateInputs(s.Config.InputFile, s.Config.ParseFormat)
	if err != nil {
		log.Fatal(err)
	}

	writer := s.createWriter()
	err = writer.Open()
	if err != nil {
		log.Fatal(err)
	}

	for index, strategy := range strategies {
		history := s.runStrategy(strategy, dateInputs)
		summary := app.CreateSummary(history)

		if index == 0 {
			err = writer.WriteSummaryHeader(summary)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = writer.WriteSummaryRow(summary)
		if err != nil {
			log.Fatal(err)
		}

		bar.Increment()
	}

	bar.Finish()

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}
