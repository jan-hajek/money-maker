package smooth

import (
	"errors"
	"fmt"
	"github.com/jelito/money-maker/app/float"
)

func NewEma(startLimit int) *EmaService {
	return &EmaService{
		startLimit: startLimit,
	}
}

type EmaService struct {
	values         []float.Float
	lastValue      float.Float
	lastValueValid bool
	startLimit     int
}

func (s *EmaService) AddStartingValue(current float.Float) error {
	if s.lastValueValid == false {
		s.values = append(s.values, current)
		return nil
	}

	return errors.New("smooth value is counting only from last value")
}

func (s *EmaService) CountSmoothValue(current, alpha float.Float) (float.Float, error) {

	if s.lastValueValid {
		s.lastValue = Ema(current, s.lastValue, alpha)
		return s.lastValue, nil
	}

	s.values = append(s.values, current)

	if len(s.values) == s.startLimit {
		s.lastValue = Avg(s.values)
		s.lastValueValid = true
		return s.lastValue, nil
	}

	return s.lastValue, errors.New(
		fmt.Sprintf("not enough starting values, expected: %d, actual: %d", s.startLimit, len(s.values)),
	)
}
