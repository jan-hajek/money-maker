package float

import (
	"github.com/jelito/money-maker/src/math/round"
)

func New(value float64) Float {
	return Float{
		value: round.Round(value, 9),
	}
}

func NewFromInt(value int) Float {
	return Float{
		value: float64(value),
	}
}

type Float struct {
	value float64
}

func (s Float) Val() float64 {
	return s.value
}

func (s Float) Add(number Float) Float {
	return New(s.value + number.value)
}

func (s Float) Sub(number Float) Float {
	return New(s.value - number.value)
}

func (s Float) Multi(number Float) Float {
	return New(s.value * number.value)
}

func (s Float) MultiFloat(number float64) Float {
	return New(s.value * number)
}

func (s Float) MultiInt(number int) Float {
	return New(s.value * float64(number))
}

func (s Float) Div(number Float) Float {
	return New(s.value / number.value)
}

func (s Float) DivInt(number int) Float {
	return New(s.value / float64(number))
}

func (s Float) DivFloat(number float64) Float {
	return New(s.value / number)
}
