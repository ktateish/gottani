# 08-cgo

This example demonstrates populating cgo functions.

Currentry it has a few restrictions:

1. Cgo file must have at most one `import "C"`.  It is permitted to have multiple
   `import "C"` in a single cgo source, but Gottani will fail to populate such files.
2. Contents of all .c files in the package are populated if cgo functions in that
   package are reachable from the entry point, `main()`.  This is done because
   we don't parse C code for now.
