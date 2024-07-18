package demo

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Set up a sample managed service",
	}
	return cmd
}
