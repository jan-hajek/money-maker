package adxEma

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/app/smooth"
	"math"
)

func New(name string, period int, alpha float.Float) *Service {
	return &Service{
		name, period, alpha,
		smooth.NewEma(period),
		smooth.NewEma(period),
		smooth.NewEma(period),
		smooth.NewEma(period),
	}
}

type Service struct {
	name            string
	period          int
	alpha           float.Float
	trueRangeSmooth *smooth.EmaService
	dmPlusSmooth    *smooth.EmaService
	dmMinusSmooth   *smooth.EmaService
	diAbsSmooth     *smooth.EmaService
}

func (s Service) Calculate(current app.IndicatorInput, history *app.History) app.IndicatorResult {
	result := Result{}

	// prvni iteraci preskakuji, je tu jen kvuli cenam z predchoziho dne
	if current.Iteration == 1 {
		return result
	}
	// FIXME - jhajek
	s.alpha = float.New(1.0 / float64(s.period))

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
	emaTrueRange, _ := s.trueRangeSmooth.CountSmoothValue(trueRange, s.alpha)
	emaDmPlus, _ := s.dmPlusSmooth.CountSmoothValue(dmPlus, s.alpha)
	emaDmMinus, _ := s.dmMinusSmooth.CountSmoothValue(dmMinus, s.alpha)

	DIAbs, DIPlus, DIMinus := s.countDmDIAbs(emaTrueRange, emaDmPlus, emaDmMinus)

	result.DIPlus = DIPlus
	result.DIMinus = DIMinus

	// v 16. - 28. (1 + 14 + 13) iteraci budu ukladat DIAbs
	if current.Iteration < 2*period {
		s.diAbsSmooth.AddStartingValue(DIAbs)

		return result
	}

	emaDiAbs, _ := s.diAbsSmooth.CountSmoothValue(DIAbs, s.alpha)

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

func (s *Service) countDmDIAbs(emaTrueRange, emaDmPlus, emaDmMinus float.Float) (float.Float, float.Float, float.Float) {
	DIPlus := float.New((100.0 * emaDmPlus.Val()) / emaTrueRange.Val())
	DIMinus := float.New((100.0 * emaDmMinus.Val()) / emaTrueRange.Val())

	DIAbs := float.New(math.Abs((DIPlus.Val() - DIMinus.Val()) / (DIPlus.Val() + DIMinus.Val())))

	return DIAbs, DIPlus, DIMinus
}
