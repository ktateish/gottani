package main

import (
	"fmt"

	"example.com/lib"
)

func main() {
	// create new lib.T
	t := lib.NewT(10)
	fmt.Println(t.Value())
	t.SetValue(t.Value() * 2)
	fmt.Println(t.Value())

	// use lib.DEFAULT through local pointer
	p := lib.DEFAULT
	fmt.Println(p.Value())
	p.SetValue(p.Value() * 2)
	fmt.Println(p.Value())

	// use lib.DEFAULT directly
	fmt.Println(lib.DEFAULT.Value())
	lib.DEFAULT.SetValue(lib.DEFAULT.Value() * 2)
	fmt.Println(lib.DEFAULT.Value())
}
