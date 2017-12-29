package app

type App struct {
	Config                  Config
	StrategyFactoryRegistry *StrategyFactoryRegistry
}

// FIXME - jhajek yzanoreni
type Config struct {
	InputFile   string                                       `yaml:"input"`
	ParseFormat string                                       `yaml:"parseFormat"`
	Outputs     map[string]interface{}                       `yaml:"outputs"`
	Strategies  map[string]map[string]map[string]interface{} `yaml:"strategies"`
}

func (s App) Run() {

	// FIXME - jhajek
	var strategy Strategy

	for strategyFactoryName, StrategyValues := range s.Config.Strategies {
		strategyFactory, err := s.StrategyFactoryRegistry.GetByName(strategyFactoryName)
		if err != nil {
			panic(err)
		}

		strategy = strategyFactory.Create(strategyFactory.GetDefaultConfig(StrategyValues))
	}

	dateInputs, err := getDateInputs(s.Config.InputFile, s.Config.ParseFormat)
	if err != nil {
		panic(err)
	}

	var lastPosition *Position

	calculators := strategy.GetCalculators()

	history := &History{
		strategy:    strategy,
		calculators: calculators,
	}

	iteration := 0
	for _, dateInput := range dateInputs {

		calculatorResults := map[string]CalculatorResult{}

		iteration++

		// FIXME - jhajek go rutiny
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

		strategyResult := strategy.Resolve(StrategyInput{
			DateInput:         dateInput,
			History:           history,
			Position:          lastPosition,
			CalculatorResults: calculatorResults,
		})

		lastPosition = createPosition(strategyResult, dateInput, lastPosition)

		historyItem := &HistoryItem{
			DateInput:         dateInput,
			CalculatorResults: calculatorResults,
			StrategyResult:    strategyResult,
			Position:          lastPosition,
		}

		history.AddItem(historyItem)

		if strategyResult.Action == CLOSE {
			lastPosition = nil
		}
	}

	summary := Summary{}
	summary.FillFromHistory(history)

	writer := s.createWriter()

	err = writer.Open()
	if err != nil {
		panic(err)
	}

	err = writer.WriteHistory(history)
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
