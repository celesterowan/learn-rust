package main

import (
	"os"

	"example.com/rust-go-course/logscope"
)

func main() {
	// Supplied boundary: library work completes before os.Exit.
	os.Exit(logscope.RunCLI(os.Args[1:], os.Stdout, os.Stderr))
}
