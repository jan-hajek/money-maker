package app

func (s App) Batch() {

	var strategies []Strategy

	for strategyFactoryName, strategyFactoryConfig := range s.Config.Strategies {
		strategyFactory, err := s.StrategyFactoryRegistry.GetByName(strategyFactoryName)
		if err != nil {
			panic(err)
		}

		for _, config := range strategyFactory.GetBatchConfigs(strategyFactoryConfig) {
			strategies = append(strategies, strategyFactory.Create(config))
		}
	}

	dateInputs, err := getDateInputs(s.Config.InputFile, s.Config.ParseFormat)
	if err != nil {
		panic(err)
	}

	var summaryList []*Summary
	for _, strategy := range strategies {
		summaryList = append(summaryList, s.runStrategy(strategy, dateInputs))
	}

	writer := s.createWriter()
	err = writer.Open()
	if err != nil {
		panic(err)
	}

	err = writer.WriteSummary(summaryList)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}
}

func (s *App) runStrategy(strategy Strategy, dateInputs []DateInput) *Summary {
	var lastPosition *Position
	indicators := strategy.GetIndicators()

	history := &History{
		strategy:   strategy,
		indicators: indicators,
	}

	iteration := 0
	for _, dateInput := range dateInputs {

		indicatorResults := map[string]IndicatorResult{}

		iteration++

		for _, c := range indicators {
			input := IndicatorInput{
				Date:       dateInput.Date,
				OpenPrice:  dateInput.OpenPrice,
				HighPrice:  dateInput.HighPrice,
				LowPrice:   dateInput.LowPrice,
				ClosePrice: dateInput.ClosePrice,
				Iteration:  iteration,
			}

			indicatorResults[c.GetName()] = c.Calculate(input, history)
		}

		strategyResult := strategy.Resolve(StrategyInput{
			DateInput:        dateInput,
			History:          history,
			Position:         lastPosition,
			IndicatorResults: indicatorResults,
		})

		lastPosition = createPosition(strategyResult, dateInput, lastPosition)

		historyItem := &HistoryItem{
			DateInput:        dateInput,
			IndicatorResults: indicatorResults,
			StrategyResult:   strategyResult,
			Position:         lastPosition,
		}

		history.AddItem(historyItem)

		if strategyResult.Action == CLOSE {
			lastPosition = nil
		}
	}

	summary := Summary{}
	summary.FillFromHistory(history)

	return &summary
}
