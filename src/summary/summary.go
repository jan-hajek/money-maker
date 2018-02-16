package summary

import (
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/printValue"
)

type Summary struct {
	StrategyPrintValues []printValue.PrintValue
	CountOfPositions    int
	CountOfProfitable   int
	CountOfLossy        int
	AvgOfPositions      float.Float
	AvgOfProfit         float.Float
	AvgOfLost           float.Float
	SumOfProfitable     float.Float
	SumOfLossy          float.Float
	Profit              float.Float
	GrossProfit         float.Float
	SuccessRatio        int
}

func CreateSummary(history interfaces.History) *Summary {
	s := &Summary{}

	s.StrategyPrintValues = history.GetStrategy().GetPrintValues()

	for _, item := range history.GetAll() {
		if item.GetStrategyResult().GetAction() == interfaces.CLOSE {
			if item.GetPosition().Profit.Val()-item.GetPosition().Costs.Val() > 0.0 {
				s.CountOfProfitable++
				s.SumOfProfitable = float.New(s.SumOfProfitable.Val() + item.GetPosition().Profit.Val() - item.GetPosition().Costs.Val())
			} else {
				s.CountOfLossy++
				s.SumOfLossy = float.New(s.SumOfLossy.Val() + item.GetPosition().Profit.Val() - item.GetPosition().Costs.Val())
			}
			s.GrossProfit = float.New(s.GrossProfit.Val() + item.GetPosition().Profit.Val())
		}
	}

	s.CountOfPositions = s.CountOfProfitable + s.CountOfLossy
	s.Profit = float.New(s.SumOfProfitable.Val() + s.SumOfLossy.Val())
	s.AvgOfProfit = float.New(s.SumOfProfitable.Val() / float64(s.CountOfProfitable))
	s.AvgOfLost = float.New(s.SumOfLossy.Val() / float64(s.CountOfLossy))
	s.AvgOfPositions = float.New(s.Profit.Val() / float64(s.CountOfPositions))
	if s.CountOfPositions != 0 {
		s.SuccessRatio = int(float64(s.CountOfProfitable) / float64(s.CountOfPositions) * 100)
	}

	return s
}
