package main

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/types"
	"os"

	"github.com/ktateish/gottani"
)

func main() {
	if err := Main(os.Args[1:]); err != nil {
		var goError error
		var foundGoError bool
	outer:
		for tmp := err; tmp != nil; tmp = errors.Unwrap(tmp) {
			switch tmp.(type) {
			case types.Error, scanner.ErrorList:
				goError = tmp
				foundGoError = true
				break outer
			}
		}
		if foundGoError {
			fmt.Fprintf(os.Stderr, "%s\n", goError)
		} else {
			fmt.Fprintf(os.Stderr, "%s: Error: %s\n", os.Args[0], err)
		}
		os.Exit(1)
	}
}

func Main(args []string) error {
	var path string
	if len(args) == 0 {
		path = "."
	} else {
		path = args[0]
	}

	b, err := gottani.Combine(path, "main")
	if err != nil {
		return err
	}
	os.Stdout.Write(b)
	return nil
}
