package servicemanagement

import (
	"fmt"

	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"cloud.google.com/go/servicemanagement/apiv1/servicemanagementpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func generateConfigReportCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "generate-config-report CONFIG",
		Short: "Generate config report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return err
			}
			defer c.Close()

			newConfig := &servicemanagementpb.ConfigRef{
				Name: args[0],
			}
			b, err := proto.Marshal(newConfig)
			if err != nil {
				return err
			}

			response, err := c.GenerateConfigReport(ctx, &servicemanagementpb.GenerateConfigReportRequest{
				NewConfig: &anypb.Any{
					TypeUrl: "type.googleapis.com/google.api.servicemanagement.v1.ConfigRef",
					Value:   b,
				},
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
