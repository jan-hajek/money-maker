package adxAvg

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"math"
)

func New(name string, period int) *Service {
	return &Service{
		name, period,
		smooth.NewSma(period),
		smooth.NewSma(period),
		smooth.NewSma(period),
		smooth.NewSma(period),
	}
}

type Service struct {
	name            string
	period          int
	trueRangeSmooth *smooth.SmaService
	dmPlusSmooth    *smooth.SmaService
	dmMinusSmooth   *smooth.SmaService
	diAbsSmooth     *smooth.SmaService
}

func (s Service) Calculate(current app.IndicatorInput, history *app.History) app.IndicatorResult {
	result := Result{}

	// prvni iteraci preskakuji, je tu jen kvuli cenam z predchoziho dne
	if current.Iteration == 1 {
		return result
	}

	period := s.period

	lastDay, _ := history.GetLastItem()
	lastInput := lastDay.DateInput

	dmPlus, dmMinus := s.countDmPlusMinus(current, lastInput)
	trueRange := s.countTrueRange(current, lastInput)

	// ve 2. - 14. iteraci se jen ukladaji hodnoty
	if current.Iteration <= period {
		s.trueRangeSmooth.AddStartingValue(trueRange)
		s.dmPlusSmooth.AddStartingValue(dmPlus)
		s.dmMinusSmooth.AddStartingValue(dmMinus)

		return result
	}

	// v 15. a dal zacinam pocitat true range, dm plus a minus
	smaTrueRange, err := s.trueRangeSmooth.CountSmoothValue(trueRange)
	smaDmPlus, _ := s.dmPlusSmooth.CountSmoothValue(dmPlus)
	smaDmMinus, _ := s.dmMinusSmooth.CountSmoothValue(dmMinus)
	if err != nil {
		panic(err)
	}

	DIAbs, DIPlus, DIMinus := s.countDmDIAbs(smaTrueRange, smaDmPlus, smaDmMinus)

	result.DIPlus = DIPlus
	result.DIMinus = DIMinus

	// v 16. - 28. (1 + 14 + 13) iteraci budu ukladat DIAbs
	if current.Iteration < 2*period {
		s.diAbsSmooth.AddStartingValue(DIAbs)

		return result
	}

	emaDiAbs, err := s.diAbsSmooth.CountSmoothValue(DIAbs)
	if err != nil {
		panic(err)
	}

	result.Adx = emaDiAbs.MultiFloat(100.0)

	return result
}

func (s Service) GetName() string {
	return s.name
}

func (s *Service) countDmPlusMinus(
	current app.IndicatorInput,
	lastInput app.DateInput,
) (float.Float, float.Float) {
	dmPlus, dmMinus := float.New(0.0), float.New(0.0)

	upMove := float.New(current.HighPrice.Val() - lastInput.HighPrice.Val())
	downMove := float.New(lastInput.LowPrice.Val() - current.LowPrice.Val())

	if upMove.Val() > downMove.Val() && upMove.Val() > 0.0 {
		dmPlus = upMove
	}

	if downMove.Val() > upMove.Val() && downMove.Val() > 0.0 {
		dmMinus = downMove
	}

	return dmPlus, dmMinus
}

func (s *Service) countTrueRange(current app.IndicatorInput, lastInput app.DateInput) float.Float {
	return float.New(math.Max(
		current.HighPrice.Val()-current.LowPrice.Val(),
		math.Max(
			math.Abs(current.HighPrice.Val()-lastInput.ClosePrice.Val()),
			math.Abs(current.LowPrice.Val()-lastInput.ClosePrice.Val()),
		),
	))
}

func (s *Service) countDmDIAbs(smaTrueRange, smaDmPlus, smaDmMinus float.Float) (float.Float, float.Float, float.Float) {
	DIPlus := float.New((100.0 * smaDmPlus.Val()) / smaTrueRange.Val())
	DIMinus := float.New((100.0 * smaDmMinus.Val()) / smaTrueRange.Val())

	DIAbs := float.New(math.Abs((DIPlus.Val() - DIMinus.Val()) / (DIPlus.Val() + DIMinus.Val())))

	return DIAbs, DIPlus, DIMinus
}
