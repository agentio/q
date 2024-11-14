package logging

import (
	"encoding/json"
	"fmt"
	"log"

	logging "cloud.google.com/go/logging/apiv2"
	"cloud.google.com/go/logging/apiv2/loggingpb"
	"github.com/spf13/cobra"
)

func listExclusionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-exclusions",
		Short: "List exclusions with the Cloud Logging API",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := logging.NewConfigClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()
			parent := args[0]
			iter := c.ListExclusions(ctx, &loggingpb.ListExclusionsRequest{
				Parent: parent,
			})
			for {
				entry, err := iter.Next()
				if err != nil {
					break
				}
				b, err := json.MarshalIndent(entry, "", "  ")
				if err != nil {
					log.Printf("%s", err)
					break
				}
				fmt.Printf("%s\n", string(b))
			}
			return nil
		},
	}
	return cmd
}
