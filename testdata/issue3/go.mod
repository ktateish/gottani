module github.com/ktateish/gottani/testdata/issue3

go 1.14

replace (
	example.com/libx => ./libx
	example.com/liby => ./liby
)

require (
	example.com/libx v0.0.0-00010101000000-000000000000
	example.com/liby v0.0.0-00010101000000-000000000000
)
