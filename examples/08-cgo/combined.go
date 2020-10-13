package main

/*
#include <math.h>

long long fact(long long n) {
	if (n <= 1) {
		return n;
	}
	return n * fact(n - 1);
}

double pi() { return M_PI; }

*/
import "C"

/*
  #cgo LDFLAGS: -lm
  long long fact(long long n);
  double pi();
*/
import "C"

import "fmt"

//line example.com/lib/lib.go:10
func Fact(n int) int {
	return int(C.fact(C.longlong(n)))
}

func Pi() float64 {
	return float64(C.pi())
}

//line main.go:9
func main() {
	for _, i := range []int{1, 2, 3, 4, 5, 10} {
		fmt.Println(Fact(i))
	}
	fmt.Println(Pi())
}
