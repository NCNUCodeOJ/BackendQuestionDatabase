package pkg

// Max returns the larger of x or y.
func Max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}
