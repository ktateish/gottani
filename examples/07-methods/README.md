# 07-methods

This example demonstrates populating methods.

All methods for a type will be populated if the type is reachable from the
entry point, `main()`.  This is done because when the type is passed to a
function that need a certain interface, it is little hard to trace which
methods are used.
