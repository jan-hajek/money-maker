package round

import m "math"

func Round(val float64, places int) (newVal float64) {
	var round float64
	pow := m.Pow(10, float64(places))
	digit := pow * val
	_, div := m.Modf(digit)
	if div >= .5 {
		round = m.Ceil(digit)
	} else {
		round = m.Floor(digit)
	}
	newVal = round / pow
	return
}
