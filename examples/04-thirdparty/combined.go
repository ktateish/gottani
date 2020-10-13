package main

import "fmt"

//line example.com/lib/lib.go:24
// Egcd(a, b) returns d, x, y:
//   d is Gcd(a,b)
//   x, y are  integers that satisfy ax + by = d
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

//line main.go:9
func main() {
	fmt.Println(Gcd(60, 24))
}
