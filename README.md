# Gottani - a tool for combining Go source code into a single .go file

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
