package lib

type foo int

type Foos []foo

func (f Foos) String() string { return "x" }
