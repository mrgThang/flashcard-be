package helpers

import "math"

func ExecuteSm2Algo(q int32, ef float32, n int32, i int32) (int32, float32, int32, int32) {
	ef = ef + (0.1 - (5.0-float32(q))*(0.08+(5.0-float32(q))*0.02))
	if ef < 1.3 {
		ef = 1.3
	}

	if q < 3 {
		n = 0
		i = 1
	} else {
		n += 1
		if n == 1 {
			i = 1
		} else if n == 2 {
			i = 6
		} else {
			i = int32(math.Round(float64(i) * float64(ef)))
			if i < 1 {
				i = 1
			}
		}
	}

	return q, ef, n, i
}
