package main

import (
	"os"

	"github.com/akashsinghal/csscgen/cmd/csscgen/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		os.Exit(1)
	}
}
