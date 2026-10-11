// create_mod builds a Go module zip from a directory using golang.org/x/mod/zip,
// the same code path the Go module proxy uses.
//
// Usage: create_mod <module-path> <version> <dir> <out.zip>
//
// This file lives under .github/ so the go tool ignores it; the seal-build
// workflow copies it into a scratch module outside the checkout to build it.
package main

import (
	"fmt"
	"os"

	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "usage: create_mod <module-path> <version> <dir> <out.zip>")
		os.Exit(2)
	}
	mv := module.Version{Path: os.Args[1], Version: os.Args[2]}
	dir := os.Args[3]
	out := os.Args[4]

	f, err := os.Create(out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create:", err)
		os.Exit(1)
	}
	if err := zip.CreateFromDir(f, mv, dir); err != nil {
		f.Close()
		fmt.Fprintln(os.Stderr, "zip.CreateFromDir:", err)
		os.Exit(1)
	}
	if err := f.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "close:", err)
		os.Exit(1)
	}
}
