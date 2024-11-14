package logging

import (
	"fmt"

	logging "cloud.google.com/go/logging/apiv2"
	"cloud.google.com/go/logging/apiv2/loggingpb"

	"github.com/spf13/cobra"
)

func listLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-logs PARENT",
		Short: "List logs with the Cloud Logging API",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := logging.NewClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()
			parent := args[0]
			iter := c.ListLogs(ctx, &loggingpb.ListLogsRequest{
				Parent: parent,
			})
			for {
				entry, err := iter.Next()
				if err != nil {
					break
				}
				fmt.Printf("%s\n", entry)
			}
			return nil
		},
	}
	return cmd
}
