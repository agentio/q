package serviceusage

import (
	"fmt"

	serviceusage "cloud.google.com/go/serviceusage/apiv1"
	"cloud.google.com/go/serviceusage/apiv1/serviceusagepb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func batchEnableServicesCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "batch-enable-services",
		Short: "Batch enable services",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := serviceusage.NewClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			response, err := c.BatchEnableServices(ctx, &serviceusagepb.BatchEnableServicesRequest{
				Parent:     args[0],
				ServiceIds: args[1:],
			})
			if err != nil {
				return err
			}
			if format == "json" {
				metadata, err := response.Metadata()
				if err != nil {
					return err
				}
				b, err := protojson.Marshal(metadata)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s", string(b))
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
