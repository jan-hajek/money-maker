package adxEma

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"math"
)

type Service struct {
	Name   string
	Period int
	alpha  float.Float
}

func (s Service) Calculate(current app.CalculatorInput, history *app.History) app.CalculatorResult {
	// prvni iteraci preskakuji, je tu jen kvuli cenam z predchoziho dne
	if current.Iteration == 1 {
		return Result{}
	}

	s.alpha = float.New(1 / float64(s.Period))

	period := s.Period

	lastDay, _ := history.GetLastItem()
	lastInput := lastDay.DateInput

	trueRange, dmPlus, dmMinus := s.countDmPlusDmMinusTrueRange(current, lastInput)

	// ve 2. - 14. iteraci se jen ukladaji hodnoty
	if current.Iteration <= period {
		return Result{
			TrueRange: trueRange,
			DmPlus:    dmPlus,
			DmMinus:   dmMinus,
		}
	}

	lastPeriodItems := history.GetLastItems(period - 1)

	// ve 15. iteraci spocitam prumery za poslednich 13 dni + aktualni
	if current.Iteration == period+1 {
		avgTrueRange, avgDmPlus, avgDmMinus := s.countAvgTrueRangeDmPlusDmMinus(
			trueRange,
			dmPlus,
			dmMinus,
			lastPeriodItems,
		)

		return Result{
			EmaTrueRange: avgTrueRange,
			EmaDmPlus:    avgDmPlus,
			EmaDmMinus:   avgDmMinus,
		}
	}

	lastValues := lastDay.CalculatorResult(s).(Result)

	emaTrueRange, emaDmPlus, emaDmMinus := s.countAvgTrueRangeSmmaDmDIAbs(
		trueRange,
		dmPlus,
		dmMinus,
		lastValues,
	)

	DIAbs, DIPlus, DIMinus := s.countDmDIAbs(emaTrueRange, emaDmPlus, emaDmMinus)

	// v 16. - 28. (1 + 14 + 13) iteraci budu ukladat DIAbs
	if current.Iteration <= 2*period {
		return Result{
			EmaTrueRange: emaTrueRange,
			EmaDmPlus:    emaDmPlus,
			EmaDmMinus:   emaDmMinus,
			DIAbs:        DIAbs,
			DIPlus:       DIPlus,
			DIMinus:      DIMinus,
		}
	}

	var smmaDIAbs float.Float

	// ve 29. (1 + 14 + 14) iteraci spocitam smmaDIAbs a ADX
	if current.Iteration == 2*period+1 {
		smmaDIAbs = s.countAvgDIAbs(DIAbs, lastPeriodItems)
		// dalsi iterace uz se pocitaji prirustkem
	} else {
		smmaDIAbs = smooth.Ema(DIAbs, lastValues.SmmaDIAbs, s.alpha)
	}

	adx := smmaDIAbs.Multi(float.New(100.0))

	return Result{
		Adx:          adx,
		EmaTrueRange: emaTrueRange,
		EmaDmPlus:    emaDmPlus,
		EmaDmMinus:   emaDmMinus,
		SmmaDIAbs:    smmaDIAbs,
		DIPlus:       DIPlus,
		DIMinus:      DIMinus,
	}
}

func (s Service) GetName() string {
	return s.Name
}

func (s *Service) countDmPlusDmMinusTrueRange(
	current app.CalculatorInput,
	lastInput app.DateInput,
) (float.Float, float.Float, float.Float) {
	dmPlus, dmMinus, trueRange := float.New(0.0), float.New(0.0), float.New(0.0)

	upMove := float.New(current.HighPrice.Val() - lastInput.HighPrice.Val())
	downMove := float.New(lastInput.LowPrice.Val() - current.LowPrice.Val())

	if upMove.Val() > downMove.Val() && upMove.Val() > 0.0 {
		dmPlus = upMove
	}

	if downMove.Val() > upMove.Val() && downMove.Val() > 0.0 {
		dmMinus = downMove
	}

	trueRange = float.New(math.Max(
		current.HighPrice.Val()-current.LowPrice.Val(),
		math.Max(
			math.Abs(current.HighPrice.Val()-lastInput.ClosePrice.Val()),
			math.Abs(current.LowPrice.Val()-lastInput.ClosePrice.Val()),
		),
	))

	return trueRange, dmPlus, dmMinus
}

func (s *Service) countAvgTrueRangeDmPlusDmMinus(
	currentTrueRange, currentDmPlus, currentDmMinus float.Float,
	lastPeriodsResults []*app.HistoryItem,
) (float.Float, float.Float, float.Float) {

	trueRangeList := []float.Float{currentTrueRange}
	dmPlusList := []float.Float{currentDmPlus}
	dmMinusList := []float.Float{currentDmMinus}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.CalculatorResult(s).(Result)
		trueRangeList = append(trueRangeList, values.TrueRange)
		dmPlusList = append(dmPlusList, values.DmPlus)
		dmMinusList = append(dmMinusList, values.DmMinus)
	}

	return smooth.Avg(trueRangeList), smooth.Avg(dmPlusList), smooth.Avg(dmMinusList)
}

func (s *Service) countAvgTrueRangeSmmaDmDIAbs(
	trueRange, dmPlus, dmMinus float.Float,
	lastValues Result,
) (float.Float, float.Float, float.Float) {
	emaTrueRange := smooth.Ema(trueRange, lastValues.EmaTrueRange, s.alpha)
	emaDmPlus := smooth.Ema(dmPlus, lastValues.EmaDmPlus, s.alpha)
	emaDmMinus := smooth.Ema(dmMinus, lastValues.EmaDmMinus, s.alpha)

	return emaTrueRange, emaDmPlus, emaDmMinus
}

func (s *Service) countDmDIAbs(emaTrueRange, emaDmPlus, emaDmMinus float.Float) (float.Float, float.Float, float.Float) {
	DIPlus := float.New((100.0 * emaDmPlus.Val()) / emaTrueRange.Val())
	DIMinus := float.New((100.0 * emaDmMinus.Val()) / emaTrueRange.Val())

	DIAbs := float.New(math.Abs((DIPlus.Val() - DIMinus.Val()) / (DIPlus.Val() + DIMinus.Val())))

	return DIAbs, DIPlus, DIMinus
}

func (s *Service) countAvgDIAbs(currentDIAbs float.Float, lastPeriodsResults []*app.HistoryItem) float.Float {
	DIAbsList := []float.Float{
		currentDIAbs,
	}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.CalculatorResult(s).(Result)
		DIAbsList = append(DIAbsList, values.DIAbs)
	}

	return smooth.Avg(DIAbsList)
}
