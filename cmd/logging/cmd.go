package logging

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logging",
		Short: "Write and manage log entries with the Cloud Logging API",
	}
	cmd.AddCommand(readCmd())
	return cmd
}

// https://github.com/googleapis/google-cloud-go/blob/main/logging/apiv2/logging_client.go
