package main

import (
	"fmt"

	"example.com/lib"
)

func main() {
	a := 3
	b := 4
	// The test harnes checks the combined.go exits with success.
	// So this call should not be executed.
	// It is OK because the purpose of this test is to confirm
	// that the combined.go can be built without any errors.
	if false {
		fmt.Printf("%d + %d = %d\n", a, b, lib.Add(a, b))
	}
}
