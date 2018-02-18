package adxEmaRSI

import (
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/math/smooth"
	"math"
)

type Service struct {
	Name       string
	Period     int
	Alpha      float.Float
	PeriodDIMA int
	DISDCount  float.Float
}

func (s Service) Calculate(current interfaces.IndicatorInput, history interfaces.History) interfaces.IndicatorResult {
	// prvni iteraci preskakuji, je tu jen kvuli cenam z predchoziho dne
	if current.Iteration == 1 {
		return Result{}
	}
	// FIXME - jhajek
	s.Alpha = float.New(1.0 / float64(s.Period))

	period := s.Period

	lastDay, _ := history.GetLastItem()
	lastInput := lastDay.GetDateInput()

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

	var DIAbs, DIPlus, DIMinus float.Float

	// ve 15. iteraci spocitam prumery za poslednich 13 dni + aktualni
	if current.Iteration == period+1 {
		avgTrueRange, avgDmPlus, avgDmMinus := s.countAvgTrueRangeDmPlusDmMinus(
			trueRange,
			dmPlus,
			dmMinus,
			lastPeriodItems,
		)

		DIAbs, DIPlus, DIMinus = s.countDmDIAbs(avgTrueRange, avgDmPlus, avgDmMinus)

		return Result{
			EmaTrueRange: avgTrueRange,
			EmaDmPlus:    avgDmPlus,
			EmaDmMinus:   avgDmMinus,
			DIAbs:        DIAbs,
			DIPlus:       DIPlus,
			DIMinus:      DIMinus,
		}
	}

	lastValues := lastDay.IndicatorResult(s).(Result)

	emaTrueRange, emaDmPlus, emaDmMinus := s.countAvgTrueRangeSmmaDmDIAbs(
		trueRange,
		dmPlus,
		dmMinus,
		lastValues,
	)

	DIAbs, DIPlus, DIMinus = s.countDmDIAbs(emaTrueRange, emaDmPlus, emaDmMinus)

	var smmaDIAbs, adx float.Float

	// ve 29. (1 + 14 + 14) iteraci spocitam smmaDIAbs a ADX
	if current.Iteration == 2*period {
		smmaDIAbs = s.countAvgDIAbs(DIAbs, lastPeriodItems)
		// dalsi iterace uz se pocitaji prirustkem
	}
	if current.Iteration > 2*period {
		smmaDIAbs = smooth.Ema(DIAbs, lastValues.SmmaDIAbs, s.Alpha)
	}

	if current.Iteration >= 2*period {
		adx = smmaDIAbs.Multi(float.New(100.0))
	}

	var DIdiff, DIdiffLB, DIdiffUB, DIdiffSD, DIdiff3, DIdiffMA float.Float

	// spocitej DIPlus - DIMinus
	if current.Iteration > period+1 {
		DIdiff = float.New(DIPlus.Val() - DIMinus.Val())
		if current.Iteration > period+3 {
			lastDIDiffItems := history.GetLastItems(2)
			DIdiff3 = s.countMA(DIdiff, lastDIDiffItems)
		}
	}
	if current.Iteration == period+1+s.PeriodDIMA {
		lastDIPeriodItems := history.GetLastItems(s.PeriodDIMA - 1)
		DIdiffSD = s.countStDeviation(DIdiff, lastDIPeriodItems)
		DIdiffMA = smooth.Ema(DIdiff, lastValues.DIDiff, float.New(1/float64(s.PeriodDIMA)))
		DIdiffLB = float.New(DIdiffMA.Val() - s.DISDCount.Val()*DIdiffSD.Val())
		DIdiffUB = float.New(DIdiffMA.Val() + s.DISDCount.Val()*DIdiffSD.Val())
	}
	if current.Iteration > period+1+s.PeriodDIMA {
		lastDIPeriodItems := history.GetLastItems(s.PeriodDIMA - 1)
		DIdiffSD = s.countStDeviation(DIdiff, lastDIPeriodItems)
		DIdiffMA = s.countMA(DIdiff, lastDIPeriodItems)
		DIdiffLB = float.New(DIdiffMA.Val() - s.DISDCount.Val()*DIdiffSD.Val())
		DIdiffUB = float.New(DIdiffMA.Val() + s.DISDCount.Val()*DIdiffSD.Val())
	}

	return Result{
		Adx:          adx,
		EmaTrueRange: emaTrueRange,
		EmaDmPlus:    emaDmPlus,
		EmaDmMinus:   emaDmMinus,
		SmmaDIAbs:    smmaDIAbs,
		DIPlus:       DIPlus,
		DIMinus:      DIMinus,
		DIAbs:        DIAbs,
		DIDiff:       DIdiff,
		DIDiff3:      DIdiff3,
		DIDiffLB:     DIdiffLB,
		DIDiffUB:     DIdiffUB,
		DIDiffMA:     DIdiffMA,
	}
}

func (s Service) GetName() string {
	return s.Name
}

func (s *Service) countDmPlusDmMinusTrueRange(
	current interfaces.IndicatorInput,
	lastInput dateInput.DateInput,
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
	lastPeriodsResults []interfaces.HistoryItem,
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

func (s *Service) countAvgTrueRangeSmmaDmDIAbs(
	trueRange, dmPlus, dmMinus float.Float,
	lastValues Result,
) (float.Float, float.Float, float.Float) {
	emaTrueRange := smooth.Ema(trueRange, lastValues.EmaTrueRange, s.Alpha)
	emaDmPlus := smooth.Ema(dmPlus, lastValues.EmaDmPlus, s.Alpha)
	emaDmMinus := smooth.Ema(dmMinus, lastValues.EmaDmMinus, s.Alpha)

	return emaTrueRange, emaDmPlus, emaDmMinus
}

func (s *Service) countDmDIAbs(emaTrueRange, emaDmPlus, emaDmMinus float.Float) (float.Float, float.Float, float.Float) {
	DIPlus := float.New((100.0 * emaDmPlus.Val()) / emaTrueRange.Val())
	DIMinus := float.New((100.0 * emaDmMinus.Val()) / emaTrueRange.Val())

	DIAbs := float.New(math.Abs((DIPlus.Val() - DIMinus.Val()) / (DIPlus.Val() + DIMinus.Val())))

	return DIAbs, DIPlus, DIMinus
}

func (s *Service) countAvgDIAbs(currentDIAbs float.Float, lastPeriodsResults []interfaces.HistoryItem) float.Float {
	DIAbsList := []float.Float{
		currentDIAbs,
	}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.IndicatorResult(s).(Result)
		DIAbsList = append(DIAbsList, values.DIAbs)
	}

	return smooth.Avg(DIAbsList)
}

func (s *Service) countStDeviation(
	currentDIdiff float.Float,
	lastPeriodsResults []interfaces.HistoryItem,
) float.Float {

	DIdiffList := []float.Float{currentDIdiff}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.IndicatorResult(s).(Result)
		DIdiffList = append(DIdiffList, values.DIDiff)
	}

	var avg, sd float.Float
	var sumOfSquares float64
	avg = smooth.Avg(DIdiffList)

	for _, value := range DIdiffList {
		sumOfSquares += math.Pow(value.Val()-avg.Val(), 2)
	}
	sd = float.New(math.Sqrt((1 / float64(len(DIdiffList)-1)) * sumOfSquares))
	return sd
}

func (s *Service) countMA(currentDIdiff float.Float, lastPeriodsResults []interfaces.HistoryItem) float.Float {
	DIdiffList := []float.Float{
		currentDIdiff,
	}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.IndicatorResult(s).(Result)
		DIdiffList = append(DIdiffList, values.DIDiff)
	}

	return smooth.Avg(DIdiffList)
}
