package mathh

import "math"

// Round a number to the nearest number of digits; I.E. 0 to round
// to an integer.
func Round(x float64, digits int) float64 {
	return math.Floor(x*math.Pow10(digits)+0.5) / math.Pow10(digits)
}
