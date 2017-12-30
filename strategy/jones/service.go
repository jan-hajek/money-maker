package jones

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"github.com/jelito/money-maker/indicator/adxAvg"
	"github.com/jelito/money-maker/indicator/adxEma"
)

type Service struct {
	config Config
	adxEma *adxEma.Service
	adxAvg *adxAvg.Service
}

func (s *Service) GetPrintValues() []app.PrintValue {
	return []app.PrintValue{
		{Label: "adxPeriod", Value: s.config.AdxPeriod},
		{Label: "alpha", Value: s.config.SmoothAlpha},
		{Label: "type", Value: s.config.SmoothType},
		{Label: "adxValOL", Value: s.config.OpenLowerAdxVal},
		{Label: "adxValOH", Value: s.config.OpenHigherAdxVal},
		{Label: "adxValC", Value: s.config.CloseAdxVal},
	}
}

func (s *Service) GetIndicators() []app.Indicator {
	if s.config.SmoothType == AVG {
		return []app.Indicator{
			s.adxAvg,
		}
	} else {
		return []app.Indicator{
			s.adxEma,
		}
	}
}

func (s Service) Resolve(input app.StrategyInput) app.StrategyResult {

	_, err := input.History.GetLastItem()
	if err != nil {
		return app.StrategyResult{
			Action: app.SKIP,
		}
	}

	var currentAdx, last3Adx, DIPlus, DIMinus float64

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
		adxValues := input.IndicatorResult(s.adxEma).(adxEma.Result)

		var last3AdxList []float.Float
		for _, item := range input.History.GetLastItems(3) {
			last3AdxList = append(last3AdxList, item.IndicatorResult(s.adxEma).(adxEma.Result).Adx)
		}
		last3Adx = smooth.Avg(last3AdxList).Val()

		currentAdx = adxValues.Adx.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	}

	if input.Position == nil {
		if currentAdx > float64(s.config.OpenHigherAdxVal) || (currentAdx >= float64(s.config.OpenLowerAdxVal) && currentAdx <= float64(s.config.OpenHigherAdxVal) && currentAdx > last3Adx) {
			var positionType app.PositionType
			if DIPlus > DIMinus {
				positionType = app.LONG
			}
			if DIPlus < DIMinus {
				positionType = app.SHORT
			}

			if positionType != "" {
				return app.StrategyResult{
					Action:       app.OPEN,
					PositionType: positionType,
					Amount:       float.New(100 / input.DateInput.ClosePrice.Val()),
				}
			}
		}
	} else {
		if input.Position.Type == app.LONG {
			if DIPlus < DIMinus || (currentAdx <= float64(s.config.CloseAdxVal) && currentAdx < last3Adx) {
				return app.StrategyResult{
					Action: app.CLOSE,
				}
			}
		}

		if input.Position.Type == app.SHORT {
			if DIPlus > DIMinus || (currentAdx <= float64(s.config.CloseAdxVal) && currentAdx < last3Adx) {
				return app.StrategyResult{
					Action: app.CLOSE,
				}
			}
		}
	}

	return app.StrategyResult{
		Action: app.SKIP,
	}

}
