package jones2

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"github.com/jelito/money-maker/indicator/adxAvg"
	"github.com/jelito/money-maker/indicator/adxEmaRSI"
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
		if DIPlus < DIMinus && (-currentDIdiff >= float64(s.config.DIOpenLevel) && currentDIdiff < lastDIdiff) {
			// (currentDIdiff <= currentDILB )
			// ||
			openShort = true
		}
	}

	if input.Position == nil {
		if (currentAdx > float64(s.config.OpenHigherAdx)) ||
			(currentAdx >= float64(s.config.OpenLowerAdx) && currentAdx <= float64(s.config.OpenHigherAdx) && adxGrowing) {
			var positionType app.PositionType
			if DIPlus > DIMinus &&
				(currentDIdiff >= float64(s.config.DIOpenLevel) && currentDIdiff > lastDIdiff) {
				//(currentDIdiff >= currentDIUB)
				// ||
				positionType = app.LONG
			}
			if DIPlus < DIMinus && (-currentDIdiff >= float64(s.config.DIOpenLevel) && currentDIdiff < lastDIdiff) {
				//(currentDIdiff <= currentDILB)
				// ||
				positionType = app.SHORT
				openShort = true
			}

			if positionType != "" {
				return app.StrategyResult{
					Action:       app.OPEN,
					PositionType: positionType,
					Amount:       float.New(100 / input.DateInput.ClosePrice.Val()),
					Costs:        s.config.Spread,
				}
			}
		}
	} else {
		if input.Position.Type == app.LONG {
			if !openLong &&
				(DIPlus <= DIMinus ||
					(currentAdx <= float64(s.config.CloseAdx) && lastAdx >= float64(s.config.CloseAdx)) ||
					//(currentDIdiff <= currentDILB) ||
					(currentDIdiff <= float64(s.config.DICloseLevel) && lastDIdiff > float64(s.config.DICloseLevel))) {
				return app.StrategyResult{
					Action: app.CLOSE,
					Costs:  input.Position.Costs,
				}
			} else {
				return app.StrategyResult{
					Action: app.CHANGE,
					Costs:  input.Position.Costs.Add(float.New(input.DateInput.ClosePrice.Val() * input.Position.Amount.Val() * (s.config.Swap.Val() / 8640))),
				}
			}
		}

		if input.Position.Type == app.SHORT {
			if !openShort &&
				(DIPlus >= DIMinus ||
					(currentAdx <= float64(s.config.CloseAdx) && lastAdx >= float64(s.config.CloseAdx)) ||
					//(currentDIdiff >= currentDIUB) ||
					(-currentDIdiff <= float64(s.config.DICloseLevel) && -lastDIdiff > float64(s.config.DICloseLevel))) {
				return app.StrategyResult{
					Action: app.CLOSE,
					Costs:  input.Position.Costs,
				}
			} else {
				return app.StrategyResult{
					Action: app.CHANGE,
					Costs:  input.Position.Costs.Add(float.New(input.DateInput.ClosePrice.Val() * input.Position.Amount.Val() * (s.config.Swap.Val() / 8640))),
				}
			}
		}
	}

	return app.StrategyResult{
		Action: app.SKIP,
	}

}
