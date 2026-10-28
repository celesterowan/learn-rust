package main

import (
	"context"
	"os"
	"os/signal"

	"example.com/rust-go-course/logscope"
)

func main() {
	os.Exit(run())
}

func run() int {
	// Supplied signal adapter: stop executes before the outer os.Exit.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	return logscope.RunServer(ctx, os.Args[1:], os.Stderr)
}
