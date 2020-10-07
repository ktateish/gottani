package main

import "fmt"

//line lib.go:24
// egcd(a, b) returns d, x, y:
//   d is gcd(a,b)
//   x, y are  integers that satisfy ax + by = d
func egcd(a, b int) (int, int, int) {
	if b == 0 {
		return a, 1, 0
	}
	d, x, y := egcd(b, a%b)
	return d, y, x - a/b*y
}

//line lib.go:35
func gcd(a, b int) int {
	d, _, _ := egcd(a, b)
	return d
}

//line main.go:7
func main() {
	fmt.Println(gcd(20835, 84561))
}
