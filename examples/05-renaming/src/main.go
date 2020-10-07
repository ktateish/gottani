package main

import (
	"fmt"
	"io"
	"os"

	"example.com/lib"
)

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
	fmt.Println(lib.Pi)

	// const
	fmt.Println(ConstB)
	fmt.Println(lib.ConstB)

	// var
	fmt.Println(VarX)
	fmt.Println(lib.VarX)

	// const and var
	fmt.Println(Y)
	fmt.Println(lib.Y)

	// function
	fmt.Println(Abs(-2))
	fmt.Println(lib.Abs(-1))

	// type
	mainT := T(10)
	libT := lib.T(20)

	mainT.Prn(os.Stdout)
	libT.Prn(os.Stdout)
}
