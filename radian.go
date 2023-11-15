package rc

import "math"

func DegreeToRadian(degree float64) (radian float64) {
	return degree * math.Pi / 180
}

func RadianToDegree(radian float64) (degree float64) {
	return radian * 180 / math.Pi
}
