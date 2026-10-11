package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Mmx233/CodeCli/cmd"
)

var Version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := cmd.NewRootCommand(Version).ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
