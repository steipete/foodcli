package main

import (
	"context"
	"os"

	"github.com/steipete/ordercli/internal/cli"
)

func run(args []string) int {
	ctx := context.Background()
	if err := cli.Run(ctx, args); err != nil {
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
