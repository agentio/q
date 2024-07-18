package logging

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logging",
		Short: "Write and manage log entries with the Google Cloud Logging API",
	}
	return cmd
}
