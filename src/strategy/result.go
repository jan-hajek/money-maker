package strategy

import (
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math/float"
)

type Result struct {
	Action        interfaces.StrategyAction
	PositionType  interfaces.PositionType
	Amount        float.Float
	Sl            float.Float
	Costs         float.Float
	ReportMessage string
}

func (s *Result) GetAction() interfaces.StrategyAction     { return s.Action }
func (s *Result) GetPositionType() interfaces.PositionType { return s.PositionType }
func (s *Result) GetAmount() float.Float                   { return s.Amount }
func (s *Result) GetSl() float.Float                       { return s.Sl }
func (s *Result) GetCosts() float.Float                    { return s.Costs }
func (s *Result) GetReportMessage() string                 { return s.ReportMessage }
