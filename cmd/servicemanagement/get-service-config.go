package servicemanagement

import (
	"fmt"

	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"cloud.google.com/go/servicemanagement/apiv1/servicemanagementpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func getServiceConfigCmd() *cobra.Command {
	var format string
	var full bool
	cmd := &cobra.Command{
		Use:   "get-service-config SERVICE CONFIG",
		Short: "Get service config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()
			view := servicemanagementpb.GetServiceConfigRequest_BASIC
			if full {
				view = servicemanagementpb.GetServiceConfigRequest_FULL
			}
			response, err := c.GetServiceConfig(ctx, &servicemanagementpb.GetServiceConfigRequest{
				ServiceName: args[0],
				ConfigId:    args[1],
				View:        view,
			})
			if err != nil {
				return err
			}
			if format == "json" {
				b, err := protojson.Marshal(response)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(b))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	cmd.Flags().BoolVar(&full, "full", false, "pass view=FULL to GetServiceConfig")

	return cmd
}
