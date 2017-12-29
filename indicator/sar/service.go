package sar

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"math"
)

type Service struct {
	Name      string
	MinimalAf float.Float
	MaximalAf float.Float
}

func (s Service) Calculate(current app.IndicatorInput, history *app.History) app.IndicatorResult {

	if current.Iteration == 1 {
		return Result{}
	}

	minimalAf := s.MinimalAf
	maximalAf := s.MaximalAf

	lastDay, _ := history.GetLastItem()
	lastInput := lastDay.DateInput
	lastIndicator := lastDay.IndicatorResult(s).(Result).Sar

	if current.Iteration == 2 {
		sar := lastInput.HighPrice
		ep := current.HighPrice
		af := minimalAf

		return Result{
			Sar:     sar,
			Ep:      ep,
			Af:      af,
			UpTrend: true,
		}
	}

	lastValues := lastDay.IndicatorResult(s).(Result)

	// treti den
	ep := float.New(0.0)
	af := float.New(0.0)

	sar := float.New(lastIndicator.Val() + (lastValues.Af.Val() * (lastValues.Ep.Val() - lastIndicator.Val())))
	lastUpTrend := lastValues.UpTrend

	if lastUpTrend {
		ep = float.New(math.Max(current.HighPrice.Val(), lastValues.Ep.Val()))
	} else {
		ep = float.New(math.Min(current.LowPrice.Val(), lastValues.Ep.Val()))
	}

	if ep.Val() != lastValues.Ep.Val() {
		af = float.New(math.Min(lastValues.Af.Val()+minimalAf.Val(), maximalAf.Val()))
	} else {
		af = lastValues.Af
	}

	if lastUpTrend && sar.Val() >= lastInput.LowPrice.Val() {
		sar = lastInput.LowPrice
	}

	upTrend := lastUpTrend
	if lastUpTrend && sar.Val() > current.LowPrice.Val() {
		upTrend = false
		sar = lastValues.Ep
		af = minimalAf
		ep = current.LowPrice
	}

	if lastUpTrend == false && sar.Val() < current.HighPrice.Val() {
		upTrend = true
		sar = lastValues.Ep
		af = minimalAf
		ep = current.HighPrice
	}

	return Result{
		Sar:     sar,
		Ep:      ep,
		Af:      af,
		UpTrend: upTrend,
	}
}

func (s Service) GetName() string {
	return s.Name
}
