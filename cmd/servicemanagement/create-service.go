package servicemanagement

import (
	"fmt"

	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"cloud.google.com/go/servicemanagement/apiv1/servicemanagementpb"
	"github.com/spf13/cobra"
)

func createServiceCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "create-service PROJECTID SERVICE",
		Short: "Create service",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			operation, err := c.CreateService(ctx, &servicemanagementpb.CreateServiceRequest{
				Service: &servicemanagementpb.ManagedService{
					ServiceName:       args[1],
					ProducerProjectId: args[0],
				},
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", operation.Name())
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
