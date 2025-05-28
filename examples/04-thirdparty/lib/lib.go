package lib

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// Egcd(a, b) returns d, x, y:
//
//	d is Gcd(a,b)
//	x, y are  integers that satisfy ax + by = d
func Egcd(a, b int) (int, int, int) {
	if b == 0 {
		return a, 1, 0
	}
	d, x, y := Egcd(b, a%b)
	return d, y, x - a/b*y
}

func Gcd(a, b int) int {
	d, _, _ := Egcd(a, b)
	return d
}
