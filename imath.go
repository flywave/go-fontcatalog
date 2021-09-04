package fontcatalog

func Abs(value int) int {
	if value == 0 {
		return 0
	}
	if value < 0 {
		return -value
	}
	return value
}

func Min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func Max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func Clamp(value, a, b int) int {
	if value < a {
		return a
	} else if value > b {
		return b
	}

	return value
}

func Sign(value int) int {
	if value == 0 {
		return 0
	}
	if value < 0 {
		return -1
	}
	return 1
}
