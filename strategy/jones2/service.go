package jones2

import (
	"fmt"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"github.com/jelito/money-maker/indicator/adxAvg"
	"github.com/jelito/money-maker/indicator/adxEmaRSI"
	"math"
)

type Service struct {
	config    Config
	adxEmaRSI *adxEmaRSI.Service
	adxAvg    *adxAvg.Service
}

func (s *Service) GetPrintValues() []app.PrintValue {
	return []app.PrintValue{
		{Label: "adxPeriod", Value: s.config.AdxPeriod},
		//{Label: "alpha", Value: s.config.SmoothAlpha},
		{Label: "type", Value: s.config.SmoothType},
		{Label: "OLadx", Value: s.config.OpenLowerAdx},
		{Label: "OHadx", Value: s.config.OpenHigherAdx},
		{Label: "Cadx", Value: s.config.CloseAdx},
		{Label: "DIOL", Value: s.config.DIOpenLevel},
		{Label: "DICL", Value: s.config.DICloseLevel},
		{Label: "SP", Value: s.config.StopProfit},
		//{Label: "DISDcount", Value: s.config.DISDCount},
		//{Label: "diPeriod", Value: s.config.PeriodDIMA},
	}
}

func (s *Service) GetIndicators() []app.Indicator {
	if s.config.SmoothType == AVG {
		return []app.Indicator{
			s.adxAvg,
		}
	} else {
		return []app.Indicator{
			s.adxEmaRSI,
		}
	}
}

func (s Service) Resolve(input app.StrategyInput) app.StrategyResult {

	lastItem, err := input.History.GetLastItem()
	if err != nil {
		return app.StrategyResult{
			Action: app.SKIP,
		}
	}

	var currentAdx, lastAdx, last3Adx, DIPlus, DIMinus, currentDIdiff, lastDIdiff float64

	if s.config.SmoothType == AVG {
		adxValues := input.IndicatorResult(s.adxAvg).(adxAvg.Result)

		var last3AdxList []float.Float
		for _, item := range input.History.GetLastItems(3) {
			last3AdxList = append(last3AdxList, item.IndicatorResult(s.adxAvg).(adxAvg.Result).Adx)
		}
		last3Adx = smooth.Avg(last3AdxList).Val()

		currentAdx = adxValues.Adx.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	} else {
		adxValues := input.IndicatorResult(s.adxEmaRSI).(adxEmaRSI.Result)

		var last3AdxList []float.Float
		for _, item := range input.History.GetLastItems(3) {
			last3AdxList = append(last3AdxList, item.IndicatorResult(s.adxEmaRSI).(adxEmaRSI.Result).Adx)
		}
		last3Adx = smooth.Avg(last3AdxList).Val()
		lastAdx = lastItem.IndicatorResult(s.adxEmaRSI).(adxEmaRSI.Result).Adx.Val()
		lastDIdiff = lastItem.IndicatorResult(s.adxEmaRSI).(adxEmaRSI.Result).DIDiff3.Val()

		currentAdx = adxValues.Adx.Val()
		currentDIdiff = adxValues.DIDiff3.Val()
		//currentDILB = adxValues.DIDiffLB.Val()
		//currentDIUB = adxValues.DIDiffUB.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	}

	var adxGrowing bool

	if currentAdx > last3Adx {
		adxGrowing = true
	} else {
		adxGrowing = false
	}

	openLong, openShort := false, false

	if (currentAdx > float64(s.config.OpenHigherAdx)) ||
		(currentAdx >= float64(s.config.OpenLowerAdx) && currentAdx <= float64(s.config.OpenHigherAdx) && adxGrowing) {
		if DIPlus > DIMinus && (currentDIdiff >= float64(s.config.DIOpenLevel) && currentDIdiff > lastDIdiff) {
			// (currentDIdiff >= currentDIUB )
			//
			openLong = true
		}
		if DIPlus < DIMinus && (-currentDIdiff >= float64(s.config.DIOpenLevel)) && (currentDIdiff < lastDIdiff) {
			// (currentDIdiff <= currentDILB )
			// ||
			openShort = true
		}
	}

	var sl float.Float

	if input.Position == nil {
		if (currentAdx > float64(s.config.OpenHigherAdx)) ||
			(currentAdx >= float64(s.config.OpenLowerAdx) && currentAdx <= float64(s.config.OpenHigherAdx) && adxGrowing) {
			var positionType app.PositionType
			if DIPlus > DIMinus &&
				(currentDIdiff >= float64(s.config.DIOpenLevel) && currentDIdiff > lastDIdiff) {
				//(currentDIdiff >= currentDIUB)
				// ||
				positionType = app.LONG
				sl = float.New(input.DateInput.ClosePrice.Val() - 0.20*input.DateInput.ClosePrice.Val())
			}
			if DIPlus < DIMinus && (-currentDIdiff >= float64(s.config.DIOpenLevel)) && (currentDIdiff < lastDIdiff) {
				//(currentDIdiff <= currentDILB)
				// ||
				positionType = app.SHORT
				openShort = true
				sl = float.New(input.DateInput.ClosePrice.Val() + 0.20*input.DateInput.ClosePrice.Val())
			}

			if positionType != "" {
				return app.StrategyResult{
					Action:       app.OPEN,
					PositionType: positionType,
					Amount:       float.New(100 / input.DateInput.ClosePrice.Val()),
					Costs:        s.config.Spread,
					Sl:           sl,
					ReportMessage: fmt.Sprintf(
						"%s, %.1f, %.1f, %s, %.3f",
						app.OPEN, lastDIdiff, currentDIdiff,
						positionType,
						sl,
					),
				}
			}
		}
	} else {
		if input.Position.Type == app.LONG {
			newSL := s.getNewSl(input.Position, int(s.config.StopProfit.Val()))

			if !openLong &&
				(DIPlus <= DIMinus ||
					(currentAdx <= float64(s.config.CloseAdx) && lastAdx >= float64(s.config.CloseAdx)) ||
					//(currentDIdiff <= currentDILB) ||
					(currentDIdiff <= float64(s.config.DICloseLevel) && lastDIdiff > float64(s.config.DICloseLevel)) ||
					(input.DateInput.LowPrice.Val() < input.Position.Sl.Val())) {
				return app.StrategyResult{
					Action: app.CLOSE,
					Costs:  input.Position.Costs,
					Sl:     newSL,
					ReportMessage: fmt.Sprintf(
						"%s, %.1f, %.1f",
						app.CLOSE, lastDIdiff, currentDIdiff,
					),
				}
			} else {
				if newSL != input.Position.Sl {
					return app.StrategyResult{
						Action: app.CHANGE,
						Costs:  input.Position.Costs.Add(float.New(input.DateInput.ClosePrice.Val() * input.Position.Amount.Val() * (s.config.Swap.Val() / 8640))),
						Sl:     newSL,
						ReportMessage: fmt.Sprintf(
							"%s, %.1f, %.1f, %.3f",
							app.CHANGE, lastDIdiff, currentDIdiff,
							newSL,
						),
					}
				} else {
					return app.StrategyResult{
						Action: app.SKIP,
						Costs:  input.Position.Costs.Add(float.New(input.DateInput.ClosePrice.Val() * input.Position.Amount.Val() * (s.config.Swap.Val() / 8640))),
						Sl:     newSL,
						ReportMessage: fmt.Sprintf(
							"%s, %.1f, %.1f",
							app.SKIP, lastDIdiff, currentDIdiff,
						),
					}
				}
			}
		}

		if input.Position.Type == app.SHORT {
			newSL := s.getNewSl(input.Position, int(s.config.StopProfit.Val()))

			if !openShort &&
				(DIPlus >= DIMinus ||
					(currentAdx <= float64(s.config.CloseAdx) && lastAdx >= float64(s.config.CloseAdx)) ||
					//(currentDIdiff >= currentDIUB) ||
					(-currentDIdiff <= float64(s.config.DICloseLevel) && -lastDIdiff > float64(s.config.DICloseLevel)) ||
					(input.DateInput.HighPrice.Val() > input.Position.Sl.Val())) {
				return app.StrategyResult{
					Action: app.CLOSE,
					Costs:  input.Position.Costs,
					Sl:     newSL,
					ReportMessage: fmt.Sprintf(
						"%s, %.1f, %.1f",
						app.CLOSE, lastDIdiff, currentDIdiff,
					),
				}
			} else {
				if newSL != input.Position.Sl {
					return app.StrategyResult{
						Action: app.CHANGE,
						Costs:  input.Position.Costs.Add(float.New(input.DateInput.ClosePrice.Val() * input.Position.Amount.Val() * (s.config.Swap.Val() / 8640))),
						Sl:     newSL,
						ReportMessage: fmt.Sprintf(
							"%s, %.1f, %.1f, %.3f",
							app.CHANGE, lastDIdiff, currentDIdiff,
							newSL,
						),
					}
				} else {
					return app.StrategyResult{
						Action: app.SKIP,
						Costs:  input.Position.Costs.Add(float.New(input.DateInput.ClosePrice.Val() * input.Position.Amount.Val() * (s.config.Swap.Val() / 8640))),
						Sl:     newSL,
						ReportMessage: fmt.Sprintf(
							"%s, %.1f, %.1f",
							app.SKIP, lastDIdiff, currentDIdiff,
						),
					}
				}

			}
		}
	}

	return app.StrategyResult{
		Action: app.SKIP,
		ReportMessage: fmt.Sprintf(
			"%s, %.1f, %.1f",
			app.SKIP, lastDIdiff, currentDIdiff,
		),
	}

}

func (s Service) getNewSl(position *app.Position, stopProfit int) float.Float {

	sl := position.Sl
	newSl := sl
	profit := 100 * position.PossibleProfitPercent.Val()
	profitFloored := int(math.Floor(100 * position.PossibleProfitPercent.Val()))
	openPrice := position.OpenPrice.Val()
	if profit > float64(stopProfit) {
		r := profitFloored - profitFloored%stopProfit - stopProfit
		if position.Type == app.SHORT {
			newSl = float.New(math.Min(sl.Val(), openPrice*(1.00-(float64(r)/100.0))))
		} else {
			newSl = float.New(math.Max(sl.Val(), openPrice*(1.00+(float64(r)/100.0))))
		}
	}

	return newSl
}
