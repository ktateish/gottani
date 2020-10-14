package lib

// DEAFULT is the default T
var DEFAULT *T

// DefaultValue returns the value of DEFAULT
func DefaultValue() int {
	return DEFAULT.Value()
}

func init() {
	// initial value of the DEFAULT is 128
	DEFAULT = NewT(128)
}

// NewT returns a pointer to the type T
func NewT(v int) *T {
	return &T{v}
}

// type T
type T struct {
	v int
}

// Value returns the value of t
func (t *T) Value() int {
	return t.v
}

// Value sets the value of t to the given v
func (t *T) SetValue(v int) {
	t.v = v
}
