package lib

// Add returns a + b
//go:noescape
//go:wasmimport
func Add(a, b int) int
