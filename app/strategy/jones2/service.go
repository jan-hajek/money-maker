package jones2

import (
	"fmt"
	"github.com/jelito/money-maker/app/indicator/adxAvg"
	"github.com/jelito/money-maker/app/indicator/adxEmaRSI"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/math/smooth"
	"github.com/jelito/money-maker/src/printValue"
	"github.com/jelito/money-maker/src/strategy"
	"math"
)

type Service struct {
	config    Config
	adxEmaRSI *adxEmaRSI.Service
	adxAvg    *adxAvg.Service
}

func (s *Service) GetPrintValues() []printValue.PrintValue {
	return []printValue.PrintValue{
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

func (s *Service) GetIndicators() []interfaces.Indicator {
	if s.config.SmoothType == AVG {
		return []interfaces.Indicator{
			s.adxAvg,
		}
	} else {
		return []interfaces.Indicator{
			s.adxEmaRSI,
		}
	}
}

func (s Service) Resolve(input interfaces.StrategyInput) interfaces.StrategyResult {

	lastItem, err := input.GetHistory().GetLastItem()
	if err != nil {
		return &strategy.Result{
			Action: interfaces.SKIP,
		}
	}

	var currentAdx, lastAdx, last3Adx, DIPlus, DIMinus, currentDIdiff, lastDIdiff float64

	if s.config.SmoothType == AVG {
		adxValues := input.IndicatorResult(s.adxAvg).(adxAvg.Result)

		var last3AdxList []float.Float
		for _, item := range input.GetHistory().GetLastItems(3) {
			last3AdxList = append(last3AdxList, item.IndicatorResult(s.adxAvg).(adxAvg.Result).Adx)
		}
		last3Adx = smooth.Avg(last3AdxList).Val()

		currentAdx = adxValues.Adx.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	} else {
		adxValues := input.IndicatorResult(s.adxEmaRSI).(adxEmaRSI.Result)

		var last3AdxList []float.Float
		for _, item := range input.GetHistory().GetLastItems(3) {
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

	if input.GetPosition() == nil {
		if (currentAdx > float64(s.config.OpenHigherAdx)) ||
			(currentAdx >= float64(s.config.OpenLowerAdx) && currentAdx <= float64(s.config.OpenHigherAdx) && adxGrowing) {
			var positionType interfaces.PositionType
			if DIPlus > DIMinus &&
				(currentDIdiff >= float64(s.config.DIOpenLevel) && currentDIdiff > lastDIdiff) {
				//(currentDIdiff >= currentDIUB)
				// ||
				positionType = interfaces.LONG
				sl = float.New(input.GetDateInput().ClosePrice.Val() - 0.20*input.GetDateInput().ClosePrice.Val())
			}
			if DIPlus < DIMinus && (-currentDIdiff >= float64(s.config.DIOpenLevel)) && (currentDIdiff < lastDIdiff) {
				//(currentDIdiff <= currentDILB)
				// ||
				positionType = interfaces.SHORT
				openShort = true
				sl = float.New(input.GetDateInput().ClosePrice.Val() + 0.20*input.GetDateInput().ClosePrice.Val())
			}

			if positionType != "" {
				return &strategy.Result{
					Action:       interfaces.OPEN,
					PositionType: positionType,
					Amount:       float.New(100 / input.GetDateInput().ClosePrice.Val()),
					Costs:        s.config.Spread,
					Sl:           sl,
					ReportMessage: fmt.Sprintf(
						"%s, %.1f, %.1f, %s, %.3f",
						interfaces.OPEN, lastDIdiff, currentDIdiff,
						positionType,
						sl,
					),
				}
			}
		}
	} else {
		if input.GetPosition().Type == interfaces.LONG {
			newSL := s.getNewSl(input.GetPosition(), int(s.config.StopProfit.Val()))

			if !openLong &&
				(DIPlus <= DIMinus ||
					(currentAdx <= float64(s.config.CloseAdx) && lastAdx >= float64(s.config.CloseAdx)) ||
					//(currentDIdiff <= currentDILB) ||
					(currentDIdiff <= float64(s.config.DICloseLevel) && lastDIdiff > float64(s.config.DICloseLevel)) ||
					(input.GetDateInput().LowPrice.Val() < input.GetPosition().Sl.Val())) {
				return &strategy.Result{
					Action: interfaces.CLOSE,
					Costs:  input.GetPosition().Costs,
					Sl:     newSL,
					ReportMessage: fmt.Sprintf(
						"%s, %.1f, %.1f",
						interfaces.CLOSE, lastDIdiff, currentDIdiff,
					),
				}
			} else {
				if newSL != input.GetPosition().Sl {
					return &strategy.Result{
						Action: interfaces.CHANGE,
						Costs:  input.GetPosition().Costs.Add(float.New(input.GetDateInput().ClosePrice.Val() * input.GetPosition().Amount.Val() * (s.config.Swap.Val() / 8640))),
						Sl:     newSL,
						ReportMessage: fmt.Sprintf(
							"%s, %.1f, %.1f, %.3f",
							interfaces.CHANGE, lastDIdiff, currentDIdiff,
							newSL,
						),
					}
				} else {
					return &strategy.Result{
						Action: interfaces.SKIP,
						Costs:  input.GetPosition().Costs.Add(float.New(input.GetDateInput().ClosePrice.Val() * input.GetPosition().Amount.Val() * (s.config.Swap.Val() / 8640))),
						Sl:     newSL,
						ReportMessage: fmt.Sprintf(
							"%s, %.1f, %.1f",
							interfaces.SKIP, lastDIdiff, currentDIdiff,
						),
					}
				}
			}
		}

		if input.GetPosition().Type == interfaces.SHORT {
			newSL := s.getNewSl(input.GetPosition(), int(s.config.StopProfit.Val()))

			if !openShort &&
				(DIPlus >= DIMinus ||
					(currentAdx <= float64(s.config.CloseAdx) && lastAdx >= float64(s.config.CloseAdx)) ||
					//(currentDIdiff >= currentDIUB) ||
					(-currentDIdiff <= float64(s.config.DICloseLevel) && -lastDIdiff > float64(s.config.DICloseLevel)) ||
					(input.GetDateInput().HighPrice.Val() > input.GetPosition().Sl.Val())) {
				return &strategy.Result{
					Action: interfaces.CLOSE,
					Costs:  input.GetPosition().Costs,
					Sl:     newSL,
					ReportMessage: fmt.Sprintf(
						"%s, %.1f, %.1f",
						interfaces.CLOSE, lastDIdiff, currentDIdiff,
					),
				}
			} else {
				if newSL != input.GetPosition().Sl {
					return &strategy.Result{
						Action: interfaces.CHANGE,
						Costs:  input.GetPosition().Costs.Add(float.New(input.GetDateInput().ClosePrice.Val() * input.GetPosition().Amount.Val() * (s.config.Swap.Val() / 8640))),
						Sl:     newSL,
						ReportMessage: fmt.Sprintf(
							"%s, %.1f, %.1f, %.3f",
							interfaces.CHANGE, lastDIdiff, currentDIdiff,
							newSL,
						),
					}
				} else {
					return &strategy.Result{
						Action: interfaces.SKIP,
						Costs:  input.GetPosition().Costs.Add(float.New(input.GetDateInput().ClosePrice.Val() * input.GetPosition().Amount.Val() * (s.config.Swap.Val() / 8640))),
						Sl:     newSL,
						ReportMessage: fmt.Sprintf(
							"%s, %.1f, %.1f",
							interfaces.SKIP, lastDIdiff, currentDIdiff,
						),
					}
				}

			}
		}
	}

	return &strategy.Result{
		Action: interfaces.SKIP,
		ReportMessage: fmt.Sprintf(
			"%s, %.1f, %.1f",
			interfaces.SKIP, lastDIdiff, currentDIdiff,
		),
	}

}

func (s Service) getNewSl(p *interfaces.Position, stopProfit int) float.Float {

	sl := p.Sl
	newSl := sl
	profit := 100 * p.PossibleProfitPercent.Val()
	profitFloored := int(math.Floor(100 * p.PossibleProfitPercent.Val()))
	openPrice := p.OpenPrice.Val()
	if profit > float64(stopProfit) {
		r := profitFloored - profitFloored%stopProfit - stopProfit
		if p.Type == interfaces.SHORT {
			newSl = float.New(math.Min(sl.Val(), openPrice*(1.00-(float64(r)/100.0))))
		} else {
			newSl = float.New(math.Max(sl.Val(), openPrice*(1.00+(float64(r)/100.0))))
		}
	}

	return newSl
}
