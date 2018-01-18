package Bollinger

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"math"
)

type Service struct {
	Name    string
	SDCount float.Float
	Period  int
}

func (s Service) Calculate(current app.IndicatorInput, history *app.History) app.IndicatorResult {
	// prvni iteraci preskakuji, je tu jen kvuli cenam z predchoziho dne
	if current.Iteration == 1 {
		return Result{}
	}

	lastDay, _ := history.GetLastItem()

	var typicalPrice, sd, ma, bl, bu float.Float

	// spocti typical price
	typicalPrice = float.New((current.HighPrice.Val() + current.LowPrice.Val() + current.ClosePrice.Val()) / 3)

	lastPeriodItems := history.GetLastItems(s.Period - 1)

	if current.Iteration == s.Period+1 {
		sd = s.countStDeviation(typicalPrice, lastPeriodItems)
		ma = s.countMA(typicalPrice, lastPeriodItems)
		bl = float.New(ma.Val() - s.SDCount.Val()*sd.Val())
		bu = float.New(ma.Val() + s.SDCount.Val()*sd.Val())
	}

	if current.Iteration > s.Period+1 {
		lastValues := lastDay.IndicatorResult(s).(Result)
		sd = s.countStDeviation(typicalPrice, lastPeriodItems)
		ma = smooth.Ema(typicalPrice, lastValues.TypicalPrice, float.New(1/float64(s.Period)))
		bl = float.New(ma.Val() - s.SDCount.Val()*sd.Val())
		bu = float.New(ma.Val() + s.SDCount.Val()*sd.Val())
	}

	return Result{
		TypicalPrice: typicalPrice,
		BL:           bl,
		BMA:          ma,
		BU:           bu,
	}
}

func (s Service) GetName() string {
	return s.Name
}

func (s *Service) countStDeviation(
	currentPrice float.Float,
	lastPeriodsResults []*app.HistoryItem,
) float.Float {

	PriceList := []float.Float{currentPrice}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.IndicatorResult(s).(Result)
		PriceList = append(PriceList, values.TypicalPrice)
	}

	var avg, sd float.Float
	var sumOfSquares float64
	avg = smooth.Avg(PriceList)

	for _, value := range PriceList {
		sumOfSquares += math.Pow(value.Val()-avg.Val(), 2)
	}
	sd = float.New(math.Sqrt((1 / float64(len(PriceList)-1)) * sumOfSquares))
	return sd
}

func (s *Service) countMA(currentPrice float.Float, lastPeriodsResults []*app.HistoryItem) float.Float {
	PriceList := []float.Float{
		currentPrice,
	}

	for _, lastPeriodResult := range lastPeriodsResults {
		values := lastPeriodResult.IndicatorResult(s).(Result)
		PriceList = append(PriceList, values.TypicalPrice)
	}

	return smooth.Avg(PriceList)
}
