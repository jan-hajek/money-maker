package adxAvgRSI

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"math"
)

type Service struct {
	Name      string
	Period    int
	PeriodRSI int
	PeriodDX  int
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

	lastValues := lastDay.IndicatorResult(s).(Result)

	var adxDiff, adxGain, adxLoss, RSI, DX, smoothAdxGain, smoothAdxLoss, gainToLoss float.Float
	// spocitej kladne a zaporne prirustky ADX
	if current.Iteration > 2*period+1 {
		adxDiff = float.New(adx.Val() - lastValues.Adx.Val())
		adxGain = float.New(math.Abs(math.Max(adxDiff.Val(), 0)))
		adxLoss = float.New(math.Abs(math.Min(adxDiff.Val(), 0)))
	}

	// vyhlad prirustky ADX a spocitej gain/loss
	if current.Iteration >= 2*period+1+s.PeriodRSI-1 {
		lastRSIPeriodItems := history.GetLastItems(s.PeriodRSI - 1)
		smoothAdxGain, smoothAdxLoss := s.countAvgAdxGainAdxLoss(adxGain, adxLoss, lastRSIPeriodItems)
		gainToLoss = smoothAdxGain.Div(smoothAdxLoss)
		RSI = float.New(100 - (100 / (1 + gainToLoss.Val())))
	}

	// vyhlad prirustky ADX a spocitej DX
	if current.Iteration >= 2*period+1+s.PeriodDX-1 {
		DX = smooth.Ema(adxDiff, lastValues.AdxDiff, float.New(1/float64(s.PeriodDX)))
	}

	return Result{
		Adx:           adx,
		EmaTrueRange:  avgTrueRange,
		EmaDmPlus:     avgDmPlus,
		EmaDmMinus:    avgDmMinus,
		SmmaDIAbs:     avgDIAbs,
		DIPlus:        DIPlus,
		DIMinus:       DIMinus,
		AdxDiff:       adxDiff,
		AdxGain:       adxGain,
		AdxLoss:       adxLoss,
		SmoothAdxGain: smoothAdxGain,
		SmoothAdxLoss: smoothAdxLoss,
		GainToLoss:    gainToLoss,
		DX:            DX,
		RSI:           RSI,
	}
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

func (s *Service) countAvgAdxGainAdxLoss(
	currentAdxGain, currentAdxLoss float.Float,
	lastRSIPeriodsResults []*app.HistoryItem,
) (float.Float, float.Float) {

	adxGainList := []float.Float{currentAdxGain}
	adxLossList := []float.Float{currentAdxLoss}

	for _, lastRSIPeriodResult := range lastRSIPeriodsResults {
		values := lastRSIPeriodResult.IndicatorResult(s).(Result)
		adxGainList = append(adxGainList, values.AdxGain)
		adxLossList = append(adxLossList, values.AdxLoss)
	}

	return smooth.Avg(adxGainList), smooth.Avg(adxLossList)
}

func (s *Service) countAvgAdxDiff(currentAdxDiff float.Float, lastPeriodsResults []*app.HistoryItem) float.Float {
	AdxDiffList := []float.Float{
		currentAdxDiff,
	}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.IndicatorResult(s).(Result)
		AdxDiffList = append(AdxDiffList, values.AdxDiff)
	}

	return smooth.Avg(AdxDiffList)
}
