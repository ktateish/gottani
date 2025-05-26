# Gottani - Combine Go source into a single .go file

## What is Gottani

Gottani is a tool for combining Go source code into a single .go file.

In competitive programming,  libraries are important to code solutions in a short
period of time, but almost all judge environment don't have third party libraries.
Many participants with Go need to do "Copy & Paste" between their library code
and solution code to be submitted.  This is really painful task for participants.

Gottani (ごった煮 meaning Ratatouille in Japanese) is a tool for such use cases.
At first, it scans all the files in the specified directory and imported
packages from there.  Then it copies the functions, constants, variables
and types reachable from the `main` function in the specified package and pastes
them into a single file renaming them if needed.  Finally it writes the combined
source code to stdout.

## Installation

```shell
# go get github.com/ktateish/gottani/cmd/gottani
```

## Usage


The following command writes the combined source code of the specified directory
to stdout.

```shell
# gottani path/to/directory
```

See also the `examples` directory.


## Note

### Handling of assembly-backed (“extern”) functions

Gottani now rewrites any Go function declaration that has no body
-i.e. a Go prototype whose real implementation lives in an .s file-
into a tiny panic stub so the flattened file always builds and links:

```
// before flattening
//go:noescape
func Add(a, b int) int

// after flattening by gottani
func Add(a, b int) int { panic("gottani: extern asm is not supported: lib.Add") }
```

#### What this means in practice

| Case    | Result                                                 |
|---------|--------------------------------------------------------|
| build   | Suceeds (the auto-generated stub prevents link errors) |
| runtime | Panics the moment an extern function is invoked        |

The stub does not replicate the original assembly routine; it only keeps
the linker happy.

#### How to stay safe

1. Verify the flattened output once with go run or go test.
    * If your program never touches the extern symbol, nothing changes.
2. Need the real assembly?
    * Keep that package as a normal import instead of flattening it, or
    * Vendor the relevant .s file alongside the generated code.

The panic message includes the fully-qualified symbol name (lib.Add above)
so you can quickly locate and refactor any accidental calls.
