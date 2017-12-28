package app

type Summary struct {
	ResolverPrintValues []PrintValue
	CountOfPositions    int
	CountOfProfitable   int
	CountOfLossy        int
	AvgOfPositions      float64
	AvgOfProfitable     float64
	AvgOfLossy          float64
	SumOfProfitable     float64
	SumOfLossy          float64
	Profit              float64
}

func (s *Summary) FillFromHistory(history *History) {
	s.ResolverPrintValues = history.resolver.GetPrintValues()

	for _, item := range history.GetAll() {
		if item.ResolverResult.Action == CLOSE {
			if item.Position.Profit > 0.0 {
				s.CountOfProfitable++
				s.SumOfProfitable += item.Position.Profit
			} else {
				s.CountOfLossy++
				s.SumOfLossy += item.Position.Profit
			}
		}
	}

	s.CountOfPositions = s.CountOfProfitable + s.CountOfProfitable
	s.Profit = s.SumOfProfitable + s.SumOfLossy
	s.AvgOfProfitable = s.SumOfProfitable / float64(s.CountOfProfitable)
	s.AvgOfLossy = s.SumOfLossy / float64(s.CountOfLossy)
	s.AvgOfPositions = s.Profit / float64(s.CountOfPositions)
}
