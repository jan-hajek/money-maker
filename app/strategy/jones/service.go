package jones

import (
	"github.com/jelito/money-maker/app/indicator/adxAvg"
	"github.com/jelito/money-maker/app/indicator/adxEma"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/math/smooth"
	"github.com/jelito/money-maker/src/printValue"
	"github.com/jelito/money-maker/src/strategy"
)

type Service struct {
	config Config
	adxEma *adxEma.Service
	adxAvg *adxAvg.Service
}

func (s *Service) GetPrintValues() []printValue.PrintValue {
	return []printValue.PrintValue{
		{Label: "adxPeriod", Value: s.config.AdxPeriod},
		{Label: "alpha", Value: s.config.SmoothAlpha},
		{Label: "type", Value: s.config.SmoothType},
		{Label: "adxValOL", Value: s.config.OpenLowerAdxVal},
		{Label: "adxValOH", Value: s.config.OpenHigherAdxVal},
		{Label: "adxValC", Value: s.config.CloseAdxVal},
	}
}

func (s *Service) GetIndicators() []interfaces.Indicator {
	if s.config.SmoothType == AVG {
		return []interfaces.Indicator{
			s.adxAvg,
		}
	} else {
		return []interfaces.Indicator{
			s.adxEma,
		}
	}
}

func (s Service) Resolve(input interfaces.StrategyInput) interfaces.StrategyResult {

	_, err := input.GetHistory().GetLastItem()
	if err != nil {
		return &strategy.Result{
			Action: interfaces.SKIP,
		}
	}

	var currentAdx, last3Adx, DIPlus, DIMinus float64

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
		adxValues := input.IndicatorResult(s.adxEma).(adxEma.Result)

		var last3AdxList []float.Float
		for _, item := range input.GetHistory().GetLastItems(3) {
			last3AdxList = append(last3AdxList, item.IndicatorResult(s.adxEma).(adxEma.Result).Adx)
		}
		last3Adx = smooth.Avg(last3AdxList).Val()

		currentAdx = adxValues.Adx.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	}

	if input.GetPosition() == nil {
		if currentAdx > float64(s.config.OpenHigherAdxVal) || (currentAdx >= float64(s.config.OpenLowerAdxVal) && currentAdx <= float64(s.config.OpenHigherAdxVal) && currentAdx > last3Adx) {
			var positionType interfaces.PositionType
			if DIPlus > DIMinus {
				positionType = interfaces.LONG
			}
			if DIPlus < DIMinus {
				positionType = interfaces.SHORT
			}

			if positionType != "" {
				return &strategy.Result{
					Action:       interfaces.OPEN,
					PositionType: positionType,
					Amount:       float.New(100 / input.GetDateInput().ClosePrice.Val()),
				}
			}
		}
	} else {
		if input.GetPosition().Type == interfaces.LONG {
			if DIPlus < DIMinus || (currentAdx <= float64(s.config.CloseAdxVal) && currentAdx < last3Adx) {
				return &strategy.Result{
					Action: interfaces.CLOSE,
				}
			}
		}

		if input.GetPosition().Type == interfaces.SHORT {
			if DIPlus > DIMinus || (currentAdx <= float64(s.config.CloseAdxVal) && currentAdx < last3Adx) {
				return &strategy.Result{
					Action: interfaces.CLOSE,
				}
			}
		}
	}

	return &strategy.Result{
		Action: interfaces.SKIP,
	}

}
