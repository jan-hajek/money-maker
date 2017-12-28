package calculator

import (
	"github.com/jelito/money-maker/app"
	"math"
)

type Adx struct {
	Name   string
	Period int
}

type AdxResult struct {
	Adx          float64
	TrueRange    float64
	AvgTrueRange float64
	DmPlus       float64
	SmmaDmPlus   float64
	DmMinus      float64
	SmmaDmMinus  float64
	DIAbs        float64
	SmmaDIAbs    float64
	DIPlus       float64
	DIMinus      float64
}

func (s AdxResult) Print() []app.PrintValue {
	return []app.PrintValue{
		{Label: "adx", Value: s.Adx},
		{Label: "diPlus", Value: s.DIPlus},
		{Label: "diMinus", Value: s.DIMinus},
	}
}

func (s Adx) Calculate(current app.CalculatorInput, history *app.History) app.CalculatorResult {
	// prvni iteraci preskakuji, je tu jen kvuli cenam z predchoziho dne
	if current.Iteration == 1 {
		return AdxResult{}
	}

	period := s.Period

	lastDay, _ := history.GetLastItem()
	lastInput := lastDay.DateInput

	dmPlus, dmMinus, trueRange := s.countDmPlusDmMinusTrueRange(current, lastInput)

	// ve 2. - 14. iteraci se jen ukladaji hodnoty
	if current.Iteration <= period {
		return AdxResult{
			TrueRange: trueRange,
			DmPlus:    dmPlus,
			DmMinus:   dmMinus,
		}
	}

	// ve 15. iteraci spocitam prumery za poslednich 13 dni + aktualni
	if current.Iteration == period+1 {
		avgTrueRange, smmaDmPlus, smmaDmMinus := s.countAfterFirstPeriod(trueRange, dmPlus, dmMinus, history, period)

		return AdxResult{
			AvgTrueRange: avgTrueRange,
			SmmaDmPlus:   smmaDmPlus,
			SmmaDmMinus:  smmaDmMinus,
		}
	}

	lastValues := lastDay.CalculatorResult(s.GetName()).(AdxResult)

	avgTrueRange, smmaDmPlus, smmaDmMinus, DIAbs, DIPlus, DIMinus := s.countAvgTrueRangeSmmaDmDIAbs(
		trueRange,
		dmPlus,
		dmMinus,
		period,
		lastValues,
	)

	// v 16. - 28. (1 + 14 + 13) iteraci budu ukladat DIAbs
	if current.Iteration <= 2*period {
		return AdxResult{
			AvgTrueRange: avgTrueRange,
			SmmaDmPlus:   smmaDmPlus,
			SmmaDmMinus:  smmaDmMinus,
			DIAbs:        DIAbs,
			DIPlus:       DIPlus,
			DIMinus:      DIMinus,
		}
	}

	smmaDIAbs := 0.0

	// ve 29. (1 + 14 + 14) iteraci spocitam smmaDIAbs a ADX
	if current.Iteration == 2*period+1 {
		smmaDIAbs = s.countSmmaDIAbs(DIAbs, history, period)
		// dalsi iterace uz se pocitaji prirustkem
	} else {
		smmaDIAbs = smooth(DIAbs, lastValues.SmmaDIAbs, period)
	}

	adx := 100 * smmaDIAbs

	return AdxResult{
		Adx:          adx,
		AvgTrueRange: avgTrueRange,
		SmmaDmPlus:   smmaDmPlus,
		SmmaDmMinus:  smmaDmMinus,
		SmmaDIAbs:    smmaDIAbs,
		DIPlus:       DIPlus,
		DIMinus:      DIMinus,
	}
}
func smooth(value float64, lastSmoothValue float64, period int) float64 {
	return (float64(period-1)*lastSmoothValue + value) / float64(period)
}

func (s Adx) GetName() string {
	return s.Name
}

func (s *Adx) countDmPlusDmMinusTrueRange(current app.CalculatorInput, lastInput app.DateInput) (float64, float64, float64) {
	dmPlus, dmMinus, trueRange := 0.0, 0.0, 0.0

	upMove := current.HighPrice - lastInput.HighPrice
	downMove := lastInput.LowPrice - current.LowPrice

	if upMove > downMove && upMove > 0.0 {
		dmPlus = upMove
	}

	if downMove > upMove && downMove > 0.0 {
		dmMinus = downMove
	}

	trueRange = math.Max(
		current.HighPrice-current.LowPrice,
		math.Max(
			math.Abs(current.HighPrice-lastInput.ClosePrice),
			math.Abs(current.LowPrice-lastInput.ClosePrice),
		),
	)

	return dmPlus, dmMinus, trueRange
}

func (s *Adx) countAfterFirstPeriod(
	trueRange, dmPlus, dmMinus float64,
	history *app.History,
	period int,

) (float64, float64, float64) {

	avgTrueRange, smmaDmPlus, smmaDmMinus := 0.0, 0.0, 0.0

	for _, lastPeriodResult := range history.GetLastItems(period - 1) {
		values := lastPeriodResult.CalculatorResult(s.GetName()).(AdxResult)
		avgTrueRange += values.AvgTrueRange
		smmaDmPlus += values.DmPlus
		smmaDmMinus += values.DmMinus
	}
	avgTrueRange += trueRange
	avgTrueRange /= float64(period)

	smmaDmPlus += dmPlus
	smmaDmPlus /= float64(period)

	smmaDmMinus += dmMinus
	smmaDmMinus /= float64(period)

	return trueRange, smmaDmPlus, smmaDmMinus
}

func (s *Adx) countAvgTrueRangeSmmaDmDIAbs(
	trueRange, dmPlus, dmMinus float64,
	period int,
	lastValues AdxResult,
) (float64, float64, float64, float64, float64, float64) {
	avgTrueRange := smooth(trueRange, lastValues.AvgTrueRange, period)

	smmaDmPlus := smooth(dmPlus, lastValues.SmmaDmPlus, period)
	smmaDmMinus := smooth(dmMinus, lastValues.SmmaDmMinus, period)

	DIPlus := (100.0 * smmaDmPlus) / avgTrueRange
	DIMinus := (100.0 * smmaDmMinus) / avgTrueRange

	DIAbs := math.Abs((DIPlus - DIMinus) / (DIPlus + DIMinus))

	return avgTrueRange, smmaDmPlus, smmaDmMinus, DIAbs, DIPlus, DIMinus
}

func (s *Adx) countSmmaDIAbs(DIAbs float64, history *app.History, period int) float64 {
	smmaDIAbs := 0.0

	for _, lastPeriodResult := range history.GetLastItems(period - 1) {
		values := lastPeriodResult.CalculatorResult(s.GetName()).(AdxResult)
		smmaDIAbs += values.DIAbs
	}
	smmaDIAbs += DIAbs
	smmaDIAbs /= float64(period)

	return smmaDIAbs
}
