package calculator

import (
	"github.com/jelito/money-maker/app"
	"math"
)

type Sar struct {
	Name      string
	MinimalAf float64
	MaximalAf float64
}

type SarResult struct {
	Sar     float64
	Af      float64
	Ep      float64
	UpTrend bool
}

func (s SarResult) Print() []app.PrintValue {
	return []app.PrintValue{
		{Label: "sar", Value: s.Sar},
	}
}

func (s Sar) Calculate(current app.CalculatorInput, history *app.History) app.CalculatorResult {

	if current.Iteration == 1 {
		return SarResult{}
	}

	minimalAf := s.MinimalAf
	maximalAf := s.MaximalAf

	lastDay, _ := history.GetLastItem()
	lastInput := lastDay.DateInput
	lastIndicator := lastDay.CalculatorResult(s.GetName()).(SarResult).Sar

	if current.Iteration == 2 {
		sar := lastInput.HighPrice
		ep := current.HighPrice
		af := minimalAf

		return SarResult{
			Sar:     sar,
			Ep:      ep,
			Af:      af,
			UpTrend: true,
		}
	}

	lastValues := lastDay.CalculatorResult(s.GetName()).(SarResult)

	// treti den
	ep := 0.0
	af := 0.0

	sar := lastIndicator + (lastValues.Af * (lastValues.Ep - lastIndicator))
	lastUpTrend := lastValues.UpTrend

	if lastUpTrend {
		ep = math.Max(current.HighPrice, lastValues.Ep)
	} else {
		ep = math.Min(current.LowPrice, lastValues.Ep)
	}

	if ep != lastValues.Ep {
		af = math.Min(lastValues.Af+minimalAf, maximalAf)
	} else {
		af = lastValues.Af
	}

	if lastUpTrend && sar >= lastInput.LowPrice {
		sar = lastInput.LowPrice
	}

	upTrend := lastUpTrend
	if lastUpTrend && sar > current.LowPrice {
		upTrend = false
		sar = lastValues.Ep
		af = minimalAf
		ep = current.LowPrice
	}

	if lastUpTrend == false && sar < current.HighPrice {
		upTrend = true
		sar = lastValues.Ep
		af = minimalAf
		ep = current.HighPrice
	}

	return SarResult{
		Sar:     sar,
		Ep:      ep,
		Af:      af,
		UpTrend: upTrend,
	}
}

func (s Sar) GetName() string {
	return s.Name
}
