package main

import (
	"fmt"

	"example.com/lib"
)

func main() {
	w := lib.NewT()
	fmt.Fprintf(w, "Hello world!")
	fmt.Println(w.String())
}
