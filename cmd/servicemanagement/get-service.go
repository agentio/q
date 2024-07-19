package servicemanagement

import (
	"fmt"

	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"cloud.google.com/go/servicemanagement/apiv1/servicemanagementpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func getServiceCmd() *cobra.Command {
	var project string
	var format string
	cmd := &cobra.Command{
		Use:   "get-service SERVICE",
		Short: "Get service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			response, err := c.GetService(ctx, &servicemanagementpb.GetServiceRequest{
				ServiceName: args[0],
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
	cmd.Flags().StringVarP(&project, "project", "p", "", "producer project")
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
