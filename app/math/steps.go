package math

func FloatSteps(min, max, step float64) []float64 {
	var results []float64

	for x := min; x <= max; x += step {
		x = Round(x, .5, 3)
		results = append(results, x)
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
