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

	ch := make(chan *app.Summary)
	x := make(chan int, 8)
	wait := make(chan bool)

	go func() {
		for i := 0; i < len(strategies); i++ {
			summary := <-ch
			<-x

			if i == 1 {

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

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}
