package logging

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logging",
		Short: "Write and manage log entries with the Cloud Logging API",
	}
	cmd.AddCommand(listBucketsCmd())
	cmd.AddCommand(listExclusionsCmd())
	cmd.AddCommand(listLinksCmd())
	cmd.AddCommand(listLogsCmd())
	cmd.AddCommand(listLogEntriesCmd())
	cmd.AddCommand(listLogMetricsCmd())
	cmd.AddCommand(listMonitoredResourceDescriptorsCmd())
	cmd.AddCommand(listSinksCmd())
	cmd.AddCommand(listViewsCmd())
	cmd.AddCommand(tailLogEntriesCmd())
	return cmd
}

// https://github.com/googleapis/google-cloud-go/blob/main/logging/apiv2/logging_client.go
