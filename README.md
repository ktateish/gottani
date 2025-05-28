# Gottani
_Combine Go source into a single .go file_

[![Go Reference](https://pkg.go.dev/badge/github.com/ktateish/gottani.svg)](https://pkg.go.dev/github.com/ktateish/gottani)
[![Go Report Card](https://goreportcard.com/badge/github.com/ktateish/gottani)](https://goreportcard.com/report/github.com/ktateish/gottani)
[![Go test](https://github.com/ktateish/gottani/actions/workflows/go-test.yml/badge.svg?branch=master)](https://github.com/ktateish/gottani/actions/workflows/go-test.yml)

## What is Gottani

Gottani is a tool for combining Go source code into a single .go file.

In competitive programming, helper libraries let you write solutions fast, but
judge environments are intentionally minimal: they usually ship only a few
popular third-party packages requested by the community. As a result, Go
competitors still end up copy-pasting their own utility code into every
submission. That manual step is a painful task.

The higher you climb on the leaderboard, the more you rely on a finely tuned
personal library—so the manual copy step becomes even more painful.

Gottani (ごった煮 meaning hotchpotch / mix-up in Japanese) eliminates that
friction. It flattens your entire project—including your custom libraries—into
a single, self-contained .go file that the judge can build without extra
dependencies.

First, it scans every file in the target directory and its imports. Then it
copies every function, constant, variable, and type reachable from `main()`
into a single file, renaming symbols as needed. Finally, it writes the combined
source to stdout.

## Disclaimer

Gottani is provided **as is, with no warranty of any kind, express or
implied.** By using this tool you agree that:

* **You are solely responsible** for verifying the generated code before
  submitting it to any online judge or contest system.

* The author and contributors **accept no liability** for lost rating points,
  disqualifications, corrupted submissions, or any other damages arising from
  the use—or inability to use—Gottani.

* Always test the flattened output with `go run` or your local judge emulator
  **well before the contest**.  If you find an issue, please open a GitHub
  issue or pull request, but understand that fixes may not arrive in time for
  your competition.

Proceed at your own risk and good luck on the leaderboard!

## Installation

```shell
$ go install github.com/ktateish/gottani/cmd/gottani@latest
```

## Usage

The following command reads source files of the specified directory that
contains a `main()` function (and any code reachable through its import chain),
then the command writes combined source code to stdout.

```shell
$ gottani path/to/directory
```

If you omit the argument, gottani will use current directory as the entry
point, i.e. main package.

```shell
$ gottani
```

See also the `examples` directory.


## Note

### Handling of assembly-backed (“extern”) functions

Gottani now rewrites any Go function declaration that has no body
-- i.e. a Go prototype whose real implementation lives in an `.s` file --
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
| build   | Succeeds (the auto-generated stub prevents link errors) |
| runtime | Panics the moment an extern function is invoked        |

The stub does not replicate the original assembly routine; it only keeps
the linker happy.

#### How to stay safe

1. Verify the flattened output once with go run or go test.
    * If your program never calls the extern symbol, nothing changes.
2. Need the real assembly?
    * Keep that package as a normal import, or
    * Vendor the relevant `.s` file alongside the generated code.

The panic message includes the fully-qualified symbol name (lib.Add above)
so you can quickly locate and refactor any accidental calls.
