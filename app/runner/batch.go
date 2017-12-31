package runner

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/jelito/money-maker/app"
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

	uiprogress.Start()
	bar := uiprogress.AddBar(len(strategies)).AppendCompleted().AppendElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("(%d/%d)", b.Current(), b.Total)
	})

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

		bar.Incr()
	}

	bar.AppendCompleted()

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}
