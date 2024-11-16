package c

func MinInt(a, b int) (min int) {
	if a <= b {
		return a
	}

	return b
}

func MinUint(a, b uint) (min uint) {
	if a <= b {
		return a
	}

	return b
}
