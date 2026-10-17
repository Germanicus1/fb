package main

import (
	"os"

	"github.com/Germanicus1/fb/internal/cli"
)

const version = "1.4.0"

func main() {
	if err := cli.Run(version); err != nil {
		cli.PrintError(os.Stderr, err)
		os.Exit(1)
	}
}
