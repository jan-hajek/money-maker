package samson

import (
	"github.com/jelito/money-maker/app/indicator/adxAvg"
	"github.com/jelito/money-maker/app/indicator/adxEma"
	"github.com/jelito/money-maker/app/indicator/sar"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/printValue"
	"github.com/jelito/money-maker/src/strategy"
)

type Service struct {
	config Config
	sar    *sar.Service
	adxEma *adxEma.Service
	adxAvg *adxAvg.Service
}

func (s *Service) GetPrintValues() []printValue.PrintValue {
	return []printValue.PrintValue{
		{Label: "sarMinAf", Value: s.config.SarMinimalAf},
		{Label: "sarMaxAf", Value: s.config.SarMaximalAf},
		{Label: "adxPeriod", Value: s.config.AdxPeriod},
		{Label: "type", Value: s.config.SmoothType},
	}
}

func (s *Service) GetIndicators() []interfaces.Indicator {
	if s.config.SmoothType == AVG {
		return []interfaces.Indicator{
			s.sar,
			s.adxAvg,
		}
	} else {
		return []interfaces.Indicator{
			s.sar,
			s.adxEma,
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

	var currentAdx, lastAdx, DIPlus, DIMinus float64

	sarValues := input.IndicatorResult(s.sar).(sar.Result)
	if s.config.SmoothType == AVG {
		adxValues := input.IndicatorResult(s.adxAvg).(adxAvg.Result)
		lastAdx = lastItem.IndicatorResult(s.adxAvg).(adxAvg.Result).Adx.Val()

		currentAdx = adxValues.Adx.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	} else {
		adxValues := input.IndicatorResult(s.adxEma).(adxEma.Result)
		lastAdx = lastItem.IndicatorResult(s.adxEma).(adxEma.Result).Adx.Val()

		currentAdx = adxValues.Adx.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	}

	currentSar := sarValues.Sar
	currentSarVal := currentSar.Val()
	lastSar := lastItem.IndicatorResult(s.sar).(sar.Result).Sar
	lastSarValue := lastSar.Val()

	currentPrice := input.GetDateInput().ClosePrice.Val()

	if input.GetPosition() == nil {

		if currentAdx > 25 || (currentAdx >= 20 && currentAdx <= 25 && currentAdx > lastAdx) {

			var positionType interfaces.PositionType
			if currentSarVal < currentPrice && DIPlus >= DIMinus {
				positionType = interfaces.LONG
			}
			if currentSarVal > currentPrice && DIPlus <= DIMinus {
				positionType = interfaces.SHORT
			}

			if positionType != "" {
				return &strategy.Result{
					Action:       interfaces.OPEN,
					PositionType: positionType,
					Amount:       float.New(1.0),
					Sl:           currentSar,
				}
			}

		}
	} else {
		if input.GetPosition().Type == interfaces.LONG {
			if currentSarVal < currentPrice && currentSarVal > lastSarValue {
				return &strategy.Result{
					Action: interfaces.CHANGE,
					Sl:     currentSar,
				}
			}

			if currentSarVal > currentPrice {
				return &strategy.Result{
					Action: interfaces.CLOSE,
					Sl:     lastSar,
				}
			}
		}

		if input.GetPosition().Type == interfaces.SHORT {
			if currentSarVal > currentPrice && currentSarVal < lastSarValue {
				return &strategy.Result{
					Action: interfaces.CHANGE,
					Sl:     currentSar,
				}
			}

			if currentSarVal < currentPrice {
				return &strategy.Result{
					Action: interfaces.CLOSE,
					Sl:     lastSar,
				}
			}
		}
	}

	return &strategy.Result{
		Action: interfaces.SKIP,
	}

}
