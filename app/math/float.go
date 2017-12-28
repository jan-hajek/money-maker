package math

func NewFloat(value float64) *Float {
	return &Float{
		value: Round(value, 3),
	}
}

type Float struct {
	value float64
}

func (s *Float) getValue() float64 {
	return s.value
}
