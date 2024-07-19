package serviceusage

import (
	"fmt"

	serviceusage "cloud.google.com/go/serviceusage/apiv1"
	"cloud.google.com/go/serviceusage/apiv1/serviceusagepb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func disableServiceCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "disable-service",
		Short: "Disable service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := serviceusage.NewClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			response, err := c.DisableService(ctx, &serviceusagepb.DisableServiceRequest{
				Name: args[0],
			})
			if err != nil {
				return err
			}
			if output == "json" {
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
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}
