// Code generated by Gottani; see https://github.com/ktateish/gottani/. DO NOT EDIT.
package main

import "fmt"

//line example.com/lib/add.go:3
// Add returns a + b
//

//line example.com/lib/add.go:7
func Add(a, b int) int { panic("gottani: extern function is not supported: lib.Add") }

//line main.go:9
func main() {
	a := 3
	b := 4
	// The test harnes checks the combined.go exits with success.
	// So this call should not be executed.
	// It is OK because the purpose of this test is to confirm
	// that the combined.go can be built without any errors.
	if false {
		fmt.Printf("%d + %d = %d\n", a, b, Add(a, b))
	}
}
