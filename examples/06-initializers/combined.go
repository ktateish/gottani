package main

import "fmt"

//line example.com/lib/lib.go:3
var X int

//line example.com/lib/lib.go:5
func initX() {
	X = 1234
}

//line example.com/lib/lib.go:9
func init() {
	initX()
}

//line example.com/lib/lib.go:13
var Y int

//line example.com/lib/lib.go:15
func initY() {
	Y = 5678
}

//line example.com/lib/lib.go:19
func init() {
	initY()
}

//line example.com/lib/lib.go:23
func GetX() int {
	return X
}

//line example.com/lib/lib.go:27
func GetY() int {
	return Y
}

//line main.go:9
func main() {
	fmt.Println(GetX())
	fmt.Println(GetY())
}
