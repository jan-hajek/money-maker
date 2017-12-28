package app

type App struct {
	Config                  Config
	ResolverFactoryRegistry *ResolverFactoryRegistry
}

// FIXME - jhajek zanoreni
type Config struct {
	InputFile   string                                       `yaml:"input"`
	ParseFormat string                                       `yaml:"parseFormat"`
	Outputs     map[string]interface{}                       `yaml:"outputs"`
	Resolvers   map[string]map[string]map[string]interface{} `yaml:"resolvers"`
}

func (s App) Run() {

	// FIXME - jhajek
	var resolver Resolver

	for resolverFactoryName, resolverValues := range s.Config.Resolvers {
		resolverFactory, err := s.ResolverFactoryRegistry.GetByName(resolverFactoryName)
		if err != nil {
			panic(err)
		}

		resolver = resolverFactory.Create(resolverFactory.GetDefaultConfig(resolverValues))
	}

	dateInputs, err := getDateInputs(s.Config.InputFile, s.Config.ParseFormat)
	if err != nil {
		panic(err)
	}

	writer := s.createWriter()

	var lastPosition *Position
	history := History{}

	iteration := 0
	for _, dateInput := range dateInputs {

		calculatorResults := map[string]CalculatorResult{}

		iteration++

		// FIXME - jhajek go rutiny
		for _, c := range resolver.GetCalculators() {
			input := CalculatorInput{
				Date:       dateInput.Date,
				OpenPrice:  dateInput.OpenPrice,
				HighPrice:  dateInput.HighPrice,
				LowPrice:   dateInput.LowPrice,
				ClosePrice: dateInput.ClosePrice,
				Iteration:  iteration,
			}

			calculatorResults[c.GetName()] = c.Calculate(input, &history)
		}

		resolverResult := resolver.Resolve(ResolverInput{
			DateInput:         dateInput,
			History:           &history,
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
	summary.FillFromHistory(&history)

	err = writer.Open()
	if err != nil {
		panic(err)
	}

	err = writer.WriteHistory(&history)
	if err != nil {
		panic(err)
	}

	err = writer.WriteSummary([]*Summary{&summary})
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}
}

func (s *App) createWriter() Writer {
	var outputs []WriterOutput

	if value, ok := s.Config.Outputs["stdout"]; ok == true && value == true {
		outputs = append(outputs, &StdOutWriterOutput{
			DateFormat: s.Config.ParseFormat,
		})
	}

	if value, ok := s.Config.Outputs["csv"]; ok == true && value != "" {
		outputs = append(outputs, &CsvWriterOutput{
			File:       value.(string),
			DateFormat: s.Config.ParseFormat,
		})
	}

	return Writer{
		outputs: outputs,
	}
}
