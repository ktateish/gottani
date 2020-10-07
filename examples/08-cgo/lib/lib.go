package lib

/*
  #cgo LDFLAGS: -lm
  long long fact(long long n);
  double pi();
*/
import "C"

func Fact(n int) int {
	return int(C.fact(C.longlong(n)))
}

func Pi() float64 {
	return float64(C.pi())
}
