package main

import (
	"fmt"

	"example.com/lib"
)

func main() {
	for _, i := range []int{1, 2, 3, 4, 5, 10} {
		fmt.Println(lib.Fact(i))
	}
	fmt.Println(lib.Pi())
}
