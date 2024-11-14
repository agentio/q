package monitoring

import (
	"encoding/json"
	"fmt"
	"log"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/spf13/cobra"
)

func listMetricDescriptorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-metric-descriptors PARENT",
		Short: "List metric descriptors with the Cloud Monitoring API",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := monitoring.NewMetricClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()
			project := args[0]
			iter := c.ListMetricDescriptors(ctx, &monitoringpb.ListMetricDescriptorsRequest{
				Name: project,
			})
			fmt.Printf("[\n")
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
				fmt.Printf("%s\n", string(b))
			}
			fmt.Printf("]\n")
			return nil

		},
	}
	return cmd
}
