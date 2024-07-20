package main

import (
	"os"

	"github.com/agent-kit/q/cmd"
)

func main() {
	if err := cmd.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
