# 05-renaming

This example demonstrates renaming feature of Gottani.

If the populated names of functions, consts, vars are duplicated with names in
main package, they'll be renamed to avoid conflict.  For example, when `lib`
pakcage has a function named `Foo` and `main` package also have the function
named `Foo`,  the format will be renamed to `lib_Foo`.  Vars, consts and types
also be renamed if needed.  In Go, top level functions, vars, consts and
types share the same scope, so they'll be renamed no matter what declaration
type of conflicting name is conflicting.

