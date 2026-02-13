package main

import (
	"os"

	"github.com/Germanicus1/fb/internal/cli"
)

const version = "1.2.0"

func main() {
	if err := cli.Run(version); err != nil {
		os.Exit(1)
	}
}
