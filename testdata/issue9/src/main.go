package main

import "example.com/lib"

type foo int

func main() {
	_ = foo(0)
	_ = lib.Foos{}
}
