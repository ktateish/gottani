package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

//line example.com/lib/lib.go:9
const lib_Pi = math.Pi

const (
	lib_ConstA = 1 << iota
	lib_ConstB
	lib_ConstC
)

var lib_VarX = "This is lib.VarX"

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

func (t lib_T) Prn(w io.Writer) {
	fmt.Fprintln(w, t)
}

//line main.go:11
const Pi = 3.14

const ConstB = "This is main.ConstB"

var VarX = "This is main.VarX"

var Y = "This is main.Y"

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type T int

func (t T) Prn(w io.Writer) {
	fmt.Fprintln(w, t)
}

func main() {
	// names of the function, consts, vars in lib will be renamed.

	// const
	fmt.Println(Pi)
	fmt.Println(lib_Pi)

	// const
	fmt.Println(ConstB)
	fmt.Println(lib_ConstB)

	// var
	fmt.Println(VarX)
	fmt.Println(lib_VarX)

	// const and var
	fmt.Println(Y)
	fmt.Println(lib_Y)

	// function
	fmt.Println(Abs(-2))
	fmt.Println(lib_Abs(-1))

	// type
	mainT := T(10)
	libT := lib_T(20)

	mainT.Prn(os.Stdout)
	libT.Prn(os.Stdout)
}
