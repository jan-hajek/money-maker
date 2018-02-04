package app

import (
	"github.com/jelito/money-maker/app/float"
)

type Strategy interface {
	Resolve(input StrategyInput) StrategyResult
	GetIndicators() []Indicator
	GetPrintValues() []PrintValue
}

type StrategyInput struct {
	DateInput        DateInput
	IndicatorResults map[string]IndicatorResult
	History          *History
	Position         *Position
}

func (s *StrategyInput) IndicatorResult(c Indicator) IndicatorResult {
	return s.IndicatorResults[c.GetName()]
}

type StrategyResult struct {
	Action        StrategyAction
	PositionType  PositionType
	Amount        float.Float
	Sl            float.Float
	Costs         float.Float
	ReportMessage string
}

type StrategyAction string

const (
	SKIP   StrategyAction = "skip"
	OPEN                  = "open"
	CLOSE                 = "close"
	CHANGE                = "change"
)

type StrategyFactory interface {
	GetDefaultConfig(config map[string]map[string]interface{}) StrategyFactoryConfig
	GetBatchConfigs(config map[string]map[string]interface{}) []StrategyFactoryConfig
	Create(config StrategyFactoryConfig) Strategy
}

type StrategyFactoryConfig interface {
}
