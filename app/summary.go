package app

import "github.com/jelito/money-maker/app/float"

type Summary struct {
	StrategyPrintValues []PrintValue
	CountOfPositions    int
	CountOfProfitable   int
	CountOfLossy        int
	AvgOfPositions      float.Float
	AvgOfProfit         float.Float
	AvgOfLost           float.Float
	SumOfProfitable     float.Float
	SumOfLossy          float.Float
	Profit              float.Float
	SuccessRatio        int
}

func (s *Summary) FillFromHistory(history *History) {
	s.StrategyPrintValues = history.strategy.GetPrintValues()

	for _, item := range history.GetAll() {
		if item.StrategyResult.Action == CLOSE {
			if item.Position.Profit.Val() > 0.0 {
				s.CountOfProfitable++
				s.SumOfProfitable = float.New(s.SumOfProfitable.Val() + item.Position.Profit.Val())
			} else {
				s.CountOfLossy++
				s.SumOfLossy = float.New(s.SumOfLossy.Val() + item.Position.Profit.Val())
			}
		}
	}

	s.CountOfPositions = s.CountOfProfitable + s.CountOfLossy
	s.Profit = float.New(s.SumOfProfitable.Val() + s.SumOfLossy.Val())
	s.AvgOfProfit = float.New(s.SumOfProfitable.Val() / float64(s.CountOfProfitable))
	s.AvgOfLost = float.New(s.SumOfLossy.Val() / float64(s.CountOfLossy))
	s.AvgOfPositions = float.New(s.Profit.Val() / float64(s.CountOfPositions))
	s.SuccessRatio = int(float64(s.CountOfProfitable) / float64(s.CountOfPositions) * 100)
}
