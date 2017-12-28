package samson

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/calculator"
)

type Service struct {
	config Config
	sar    *calculator.Sar
	adx    *calculator.Adx
}

func (s *Service) GetPrintValues() []app.PrintValue {
	return []app.PrintValue{
		{Label: "sarMinAf", Value: s.config.SarMinimalAf},
		{Label: "sarMaxAf", Value: s.config.SarMaximalAf},
		{Label: "adxPeriod", Value: s.config.AdxPeriod},
	}
}

func (s *Service) GetCalculators() []app.Calculator {
	return []app.Calculator{
		s.sar,
		s.adx,
	}
}

func (s Service) Resolve(input app.ResolverInput) app.ResolverResult {

	sarValues := input.CalculatorResult(s.sar.GetName()).(calculator.SarResult)
	adxValues := input.CalculatorResult(s.adx.GetName()).(calculator.AdxResult)

	lastItem, err := input.History.GetLastItem()
	if err != nil {
		return app.ResolverResult{
			Action: app.SKIP,
		}
	}

	currentAdx := adxValues.Adx
	lastAdx := lastItem.CalculatorResult(s.adx.GetName()).(calculator.AdxResult).Adx

	currentSar := sarValues.Sar
	lastSar := lastItem.CalculatorResult(s.sar.GetName()).(calculator.SarResult).Sar

	currentPrice := input.DateInput.ClosePrice

	if input.Position == nil {

		if currentAdx > 25 || (currentAdx >= 20 && currentAdx <= 25 && currentAdx > lastAdx) {

			var positionType app.PositionType
			if currentSar < currentPrice && adxValues.DIPlus >= adxValues.DIMinus {
				positionType = app.LONG
			}
			if currentSar > currentPrice && adxValues.DIPlus <= adxValues.DIMinus {
				positionType = app.SHORT
			}

			if positionType != "" {
				return app.ResolverResult{
					Action:       app.OPEN,
					PositionType: positionType,
					Amount:       1.0,
					Sl:           currentSar,
				}
			}

		}
	} else {
		if input.Position.Type == app.LONG {
			if currentSar < currentPrice && currentSar > lastSar {
				return app.ResolverResult{
					Action: app.CHANGE,
					Sl:     currentSar,
				}
			}

			if currentSar > currentPrice {
				return app.ResolverResult{
					Action: app.CLOSE,
					Sl:     lastSar,
				}
			}
		}

		if input.Position.Type == app.SHORT {
			if currentSar > currentPrice && currentSar < lastSar {
				return app.ResolverResult{
					Action: app.CHANGE,
					Sl:     currentSar,
				}
			}

			if currentSar < currentPrice {
				return app.ResolverResult{
					Action: app.CLOSE,
					Sl:     lastSar,
				}
			}
		}
	}

	return app.ResolverResult{
		Action: app.SKIP,
	}

}
