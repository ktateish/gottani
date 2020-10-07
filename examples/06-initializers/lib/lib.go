package lib

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
