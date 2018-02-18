package math

import (
	"github.com/jelito/money-maker/src/math/float"
)

func FloatSteps(min, max, step float.Float) []float.Float {
	var results []float.Float

	for x := min.Val(); x <= max.Val(); x += step.Val() {
		results = append(results, float.New(x))
	}

	return results
}

func IntSteps(min, max, step int) []int {
	var results []int

	for x := min; x <= max; x += step {
		results = append(results, x)
	}

	return results
}
