package main

import "fmt"

// =============================================================================
// Populated Libiraries
// =============================================================================

var (
//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:3
	X int

//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:13
	Y int

//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:13
)

//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:9
func init() {
	initX()
}

//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:19
func init() {
	initY()
}

//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:5
func initX() {
	X = 1234
}

//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:15
func initY() {
	Y = 5678
}

//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:23
func GetX() int {
	return X
}

//line github.com/ktateish/gottani/examples/06-initializers/lib/lib.go:27
func GetY() int {
	return Y
}

// =============================================================================
// Original Main Package
// =============================================================================

//line main.go:9
func main() {
	fmt.Println(GetX())
	fmt.Println(GetY())
}
