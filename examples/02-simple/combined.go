package main

import "fmt"

//line lib.go:10
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//line main.go:7
func main() {
	fmt.Println(min(1, 100))
}