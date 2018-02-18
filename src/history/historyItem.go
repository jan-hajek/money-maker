package history

import (
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/interfaces"
)

type HistoryItem struct {
	DateInput        dateInput.DateInput
	IndicatorResults map[string]interfaces.IndicatorResult
	StrategyResult   interfaces.StrategyResult
	Position         *interfaces.Position
}

func (s *HistoryItem) GetDateInput() dateInput.DateInput {
	return s.DateInput
}

func (s *HistoryItem) GetStrategyResult() interfaces.StrategyResult {
	return s.StrategyResult
}

func (s *HistoryItem) GetPosition() *interfaces.Position {
	return s.Position
}

func (s *HistoryItem) GetIndicatorResults() map[string]interfaces.IndicatorResult {
	return s.IndicatorResults
}

func (s *HistoryItem) IndicatorResult(c interfaces.Indicator) interfaces.IndicatorResult {
	item := s.IndicatorResults[c.GetName()]

	return item
}
