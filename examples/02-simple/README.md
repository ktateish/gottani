# 02-simple

The src/ directory has main.go and lib.go.  the main.go uses only one function
in the lib.go.

The difference between 01-simple/src and 02-simple/src is only a following line
and Gottani populates the proper function for each case.

```diff
diff --git a/01-simple/src/main.go b/02-simple/src/main.go
index f5ade44..eeb16be 100644
--- a/01-simple/src/main.go
+++ b/02-simple/src/main.go
@@ -5,5 +5,5 @@ import (
 )

 func main() {
-       fmt.Println(max(1, 100))
+       fmt.Println(min(1, 100))
 }
```
