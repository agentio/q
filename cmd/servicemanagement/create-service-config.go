package servicemanagement

import (
	"fmt"
	"os"

	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"cloud.google.com/go/servicemanagement/apiv1/servicemanagementpb"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/encoding/protojson"
)

func createServiceConfigCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "create-service-config SERVICE FILE",
		Short: "Create service config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()

			b, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}
			var service serviceconfig.Service
			err = protojson.Unmarshal(b, &service)
			if err != nil {
				return err
			}
			response, err := c.CreateServiceConfig(ctx, &servicemanagementpb.CreateServiceConfigRequest{
				ServiceName:   args[0],
				ServiceConfig: &service,
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
	return cmd
}
