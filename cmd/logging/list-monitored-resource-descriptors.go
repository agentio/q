package logging

import (
	"encoding/json"
	"fmt"
	"log"

	logging "cloud.google.com/go/logging/apiv2"
	"cloud.google.com/go/logging/apiv2/loggingpb"

	"github.com/spf13/cobra"
)

func listMonitoredResourceDescriptorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-monitored-resource-descriptors",
		Short: "List the descriptors for monitored resource types used by Logging",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := logging.NewClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()
			iter := c.ListMonitoredResourceDescriptors(ctx, &loggingpb.ListMonitoredResourceDescriptorsRequest{})
			fmt.Printf("[")
			for i := 0; true; i += 1 {
				entry, err := iter.Next()
				if err != nil {
					break
				}
				b, err := json.MarshalIndent(entry, "", "  ")
				if err != nil {
					log.Printf("%s", err)
					break
				}
				if i > 0 {
					fmt.Printf(",")
				}
				fmt.Printf("\n%s", string(b))
			}
			fmt.Printf("\n]\n")
			return nil
		},
	}
	return cmd
}
