package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// We need to:
	// - Setup the app
	// - Setup config
	// - Setup an LLM client
	// - Setup the read-eval-print loop (REPL)
	// -

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, config.load()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
