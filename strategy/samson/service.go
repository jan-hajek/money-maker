package samson

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/calculator/adxAvg"
	"github.com/jelito/money-maker/calculator/adxEma"
	"github.com/jelito/money-maker/calculator/sar"
)

type Service struct {
	config Config
	sar    *sar.Service
	adxEma *adxEma.Service
	adxAvg *adxAvg.Service
}

func (s *Service) GetPrintValues() []app.PrintValue {
	return []app.PrintValue{
		{Label: "sarMinAf", Value: s.config.SarMinimalAf},
		{Label: "sarMaxAf", Value: s.config.SarMaximalAf},
		{Label: "adxPeriod", Value: s.config.AdxPeriod},
		{Label: "type", Value: s.config.SmoothType},
	}
}

func (s *Service) GetCalculators() []app.Calculator {
	if s.config.SmoothType == AVG {
		return []app.Calculator{
			s.sar,
			s.adxAvg,
		}
	} else {
		return []app.Calculator{
			s.sar,
			s.adxEma,
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

	var currentAdx, lastAdx, DIPlus, DIMinus float64

	sarValues := input.CalculatorResult(s.sar.GetName()).(sar.Result)
	if s.config.SmoothType == AVG {
		adxValues := input.CalculatorResult(s.adxAvg.GetName()).(adxAvg.Result)
		lastAdx = lastItem.CalculatorResult(s.adxAvg.GetName()).(adxAvg.Result).Adx.Val()

		currentAdx = adxValues.Adx.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	} else {
		adxValues := input.CalculatorResult(s.adxEma.GetName()).(adxEma.Result)
		lastAdx = lastItem.CalculatorResult(s.adxEma.GetName()).(adxEma.Result).Adx.Val()

		currentAdx = adxValues.Adx.Val()
		DIPlus = adxValues.DIPlus.Val()
		DIMinus = adxValues.DIMinus.Val()
	}

	currentSar := sarValues.Sar
	currentSarVal := currentSar.Val()
	lastSar := lastItem.CalculatorResult(s.sar.GetName()).(sar.Result).Sar
	lastSarValue := lastSar.Val()

	currentPrice := input.DateInput.ClosePrice.Val()

	if input.Position == nil {

		if currentAdx > 25 || (currentAdx >= 20 && currentAdx <= 25 && currentAdx > lastAdx) {

			var positionType app.PositionType
			if currentSarVal < currentPrice && DIPlus >= DIMinus {
				positionType = app.LONG
			}
			if currentSarVal > currentPrice && DIPlus <= DIMinus {
				positionType = app.SHORT
			}

			if positionType != "" {
				return app.StrategyResult{
					Action:       app.OPEN,
					PositionType: positionType,
					Amount:       float.New(1.0),
					Sl:           currentSar,
				}
			}

		}
	} else {
		if input.Position.Type == app.LONG {
			if currentSarVal < currentPrice && currentSarVal > lastSarValue {
				return app.StrategyResult{
					Action: app.CHANGE,
					Sl:     currentSar,
				}
			}

			if currentSarVal > currentPrice {
				return app.StrategyResult{
					Action: app.CLOSE,
					Sl:     lastSar,
				}
			}
		}

		if input.Position.Type == app.SHORT {
			if currentSarVal > currentPrice && currentSarVal < lastSarValue {
				return app.StrategyResult{
					Action: app.CHANGE,
					Sl:     currentSar,
				}
			}

			if currentSarVal < currentPrice {
				return app.StrategyResult{
					Action: app.CLOSE,
					Sl:     lastSar,
				}
			}
		}
	}

	return app.StrategyResult{
		Action: app.SKIP,
	}

}
