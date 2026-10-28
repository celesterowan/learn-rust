package main

import (
	"fmt"
	"os"
)

func main() {
	// TODO G5: Trace this supplied wiring while implementing Run in cli.go.
	output := commandResult(Run(os.Args[1:]))
	fmt.Fprint(os.Stdout, output.stdout)
	fmt.Fprint(os.Stderr, output.stderr)
	os.Exit(output.status)
}
