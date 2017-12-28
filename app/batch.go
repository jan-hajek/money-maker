package app

func (s App) Batch() {

	var resolvers []Resolver

	for resolverFactoryName, resolverFactoryConfig := range s.Config.Resolvers {
		resolverFactory, err := s.ResolverFactoryRegistry.GetByName(resolverFactoryName)
		if err != nil {
			panic(err)
		}

		for _, config := range resolverFactory.GetBatchConfigs(resolverFactoryConfig) {
			resolvers = append(resolvers, resolverFactory.Create(config))
		}
	}

	dateInputs, err := getDateInputs(s.Config.InputFile, s.Config.ParseFormat)
	if err != nil {
		panic(err)
	}

	var summaryList []*Summary
	for _, resolver := range resolvers {
		summaryList = append(summaryList, s.runResolver(resolver, dateInputs))
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

func (s *App) runResolver(resolver Resolver, dateInputs []DateInput) *Summary {
	var lastPosition *Position
	calculators := resolver.GetCalculators()

	history := &History{
		resolver:    resolver,
		calculators: calculators,
	}

	iteration := 0
	for _, dateInput := range dateInputs {

		calculatorResults := map[string]CalculatorResult{}

		iteration++

		for _, c := range calculators {
			input := CalculatorInput{
				Date:       dateInput.Date,
				OpenPrice:  dateInput.OpenPrice,
				HighPrice:  dateInput.HighPrice,
				LowPrice:   dateInput.LowPrice,
				ClosePrice: dateInput.ClosePrice,
				Iteration:  iteration,
			}

			calculatorResults[c.GetName()] = c.Calculate(input, history)
		}

		resolverResult := resolver.Resolve(ResolverInput{
			DateInput:         dateInput,
			History:           history,
			Position:          lastPosition,
			CalculatorResults: calculatorResults,
		})

		lastPosition = createPosition(resolverResult, dateInput, lastPosition)

		historyItem := &HistoryItem{
			DateInput:         dateInput,
			CalculatorResults: calculatorResults,
			ResolverResult:    resolverResult,
			Position:          lastPosition,
		}

		history.AddItem(historyItem)

		if resolverResult.Action == CLOSE {
			lastPosition = nil
		}
	}

	summary := Summary{}
	summary.FillFromHistory(history)

	return &summary
}
