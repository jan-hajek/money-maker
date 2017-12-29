package adxAvg

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"math"
)

type Service struct {
	Name   string
	Period int
}

func (s Service) Calculate(current app.IndicatorInput, history *app.History) app.IndicatorResult {
	result := Result{}

	// prvni iteraci preskakuji, je tu jen kvuli cenam z predchoziho dne
	if current.Iteration == 1 {
		return result
	}

	period := s.Period

	lastDay, _ := history.GetLastItem()
	lastInput := lastDay.DateInput

	trueRange, dmPlus, dmMinus := s.countDmPlusDmMinusTrueRange(current, lastInput)

	result.TrueRange = trueRange
	result.DmPlus = dmPlus
	result.DmMinus = dmMinus

	// ve 2. - 14. iteraci se jen ukladaji hodnoty
	if current.Iteration <= period {
		return result
	}

	lastPeriodItems := history.GetLastItems(period - 1)

	// od 15. iteraci spocitam prumery za poslednich 13 dni + aktualni
	avgTrueRange, avgDmPlus, avgDmMinus := s.countAvgTrueRangeDmPlusDmMinus(trueRange, dmPlus, dmMinus, lastPeriodItems)

	DIAbs, DIPlus, DIMinus := s.countDmDIAbs(avgTrueRange, avgDmPlus, avgDmMinus)

	result.DIAbs = DIAbs
	result.DIPlus = DIPlus
	result.DIMinus = DIMinus

	// v 16. - 28. (1 + 14 + 13) iteraci budu ukladat DIAbs
	if current.Iteration <= 2*period {
		return result
	}

	// od 29. (1 + 14 + 14) iteraci spocitam smmaDIAbs a ADX
	avgDIAbs := s.countAvgDIAbs(DIAbs, lastPeriodItems)

	adx := avgDIAbs.MultiFloat(100.0)

	result.Adx = adx

	return result
}

func (s Service) GetName() string {
	return s.Name
}

func (s *Service) countDmPlusDmMinusTrueRange(
	current app.IndicatorInput,
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
		values := lastPeriodResult.IndicatorResult(s).(Result)
		trueRangeList = append(trueRangeList, values.TrueRange)
		dmPlusList = append(dmPlusList, values.DmPlus)
		dmMinusList = append(dmMinusList, values.DmMinus)
	}

	return smooth.Avg(trueRangeList), smooth.Avg(dmPlusList), smooth.Avg(dmMinusList)
}

func (s *Service) countDmDIAbs(avgTrueRange, avgDmPlus, avgDmMinus float.Float) (float.Float, float.Float, float.Float) {
	DIPlus := float.New((100.0 * avgDmPlus.Val()) / avgTrueRange.Val())
	DIMinus := float.New((100.0 * avgDmMinus.Val()) / avgTrueRange.Val())

	DIAbs := float.New(math.Abs((DIPlus.Val() - DIMinus.Val()) / (DIPlus.Val() + DIMinus.Val())))

	return DIAbs, DIPlus, DIMinus
}

func (s *Service) countAvgDIAbs(currentDIAbs float.Float, lastPeriodsResults []*app.HistoryItem) float.Float {
	DIAbsList := []float.Float{
		currentDIAbs,
	}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.IndicatorResult(s).(Result)
		DIAbsList = append(DIAbsList, values.DIAbs)
	}

	return smooth.Avg(DIAbsList)
}
