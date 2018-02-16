package smooth

import "github.com/jelito/money-maker/src/math/float"

func Ema(actual, last, alpha float.Float) float.Float {
	return float.New(((1.0 - alpha.Val()) * last.Val()) + (alpha.Val() * actual.Val()))
}

func Avg(values []float.Float) float.Float {
	avg := 0.0
	for _, value := range values {
		avg += value.Val()
	}
	avg /= float64(len(values))

	return float.New(avg)
}
