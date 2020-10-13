package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

//line example.com/lib/lib.go:9
const lib_Pi = math.Pi

//line example.com/lib/lib.go:11
const (
	lib_ConstA = 1 << iota
	lib_ConstB
	lib_ConstC
)

//line example.com/lib/lib.go:17
var lib_VarX = "This is lib.VarX"

//line example.com/lib/lib.go:19
const lib_Y = 123

//line example.com/lib/lib.go:35
func lib_Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

//line example.com/lib/lib.go:60
type lib_T float64

//line example.com/lib/lib.go:62
func (t lib_T) Prn(w io.Writer) {
	fmt.Fprintln(w, t)
}

//line main.go:11
const Pi = 3.14

//line main.go:13
const ConstB = "This is main.ConstB"

//line main.go:15
var VarX = "This is main.VarX"

//line main.go:17
var Y = "This is main.Y"

//line main.go:19
func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

//line main.go:26
type T int

//line main.go:28
func (t T) Prn(w io.Writer) {
	fmt.Fprintln(w, t)
}

//line main.go:32
func main() {

//line main.go:36
	fmt.Println(Pi)
	fmt.Println(lib_Pi)

//line main.go:40
	fmt.Println(ConstB)
	fmt.Println(lib_ConstB)

//line main.go:44
	fmt.Println(VarX)
	fmt.Println(lib_VarX)

//line main.go:48
	fmt.Println(Y)
	fmt.Println(lib_Y)

//line main.go:52
	fmt.Println(Abs(-2))
	fmt.Println(lib_Abs(-1))

//line main.go:56
	mainT := T(10)
	libT := lib_T(20)

	mainT.Prn(os.Stdout)
	libT.Prn(os.Stdout)
}
