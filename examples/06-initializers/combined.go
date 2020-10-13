package main

import "fmt"

//line example.com/lib/lib.go:3
var X int

func initX() {
	X = 1234
}

func init() {
	initX()
}

var Y int

func initY() {
	Y = 5678
}

func init() {
	initY()
}

func GetX() int {
	return X
}

func GetY() int {
	return Y
}

//line main.go:9
func main() {
	fmt.Println(GetX())
	fmt.Println(GetY())
}
