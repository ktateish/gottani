package main

import "fmt"

//line lib.go:3
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

//line main.go:7
func main() {
	fmt.Println(max(1, 100))
}
