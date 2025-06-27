package lib

import "fmt"

type Type[T any] struct{}

func New[T any](v T) Type[T] {
	return Type[T]{}
}

func (t Type[T]) Foo() {
	fmt.Printf("Hello from %T\n", t)
}
