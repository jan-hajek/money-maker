package smooth

import (
	"errors"
	"fmt"
	"github.com/jelito/money-maker/app/float"
)

func NewSma(startLimit int) *SmaService {
	return &SmaService{
		startLimit: startLimit,
	}
}

type SmaService struct {
	values     []float.Float
	startLimit int
}

func (s *SmaService) AddStartingValue(current float.Float) error {
	if len(s.values) < s.startLimit {
		s.values = append(s.values, current)
		return nil
	}

	return errors.New("")
}

func (s *SmaService) CountSmoothValue(current float.Float) (float.Float, error) {

	s.values = append(s.values, current)

	if len(s.values) == s.startLimit {
		avg := Avg(s.values)

		// FIXME - jhajek leak
		s.values = s.values[1:]

		return avg, nil

	}

	return float.New(0.0), errors.New(
		fmt.Sprintf("not enough starting values, expected: %d, actual: %d", s.startLimit, len(s.values)),
	)
}
