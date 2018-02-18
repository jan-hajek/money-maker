package strategy

import (
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/interfaces"
)

type Input struct {
	DateInput        dateInput.DateInput
	IndicatorResults map[string]interfaces.IndicatorResult
	History          interfaces.History
	Position         *interfaces.Position
}

func (s *Input) GetDateInput() dateInput.DateInput { return s.DateInput }
func (s *Input) GetPosition() *interfaces.Position { return s.Position }
func (s *Input) GetHistory() interfaces.History    { return s.History }
func (s *Input) IndicatorResult(c interfaces.Indicator) interfaces.IndicatorResult {
	return s.IndicatorResults[c.GetName()]
}
