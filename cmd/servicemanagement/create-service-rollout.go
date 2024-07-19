package servicemanagement

import (
	"fmt"

	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"cloud.google.com/go/servicemanagement/apiv1/servicemanagementpb"
	"github.com/spf13/cobra"
)

func createServiceRolloutCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "create-service-rollout SERVICE CONFIG",
		Short: "Create service rollout",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			operation, err := c.CreateServiceRollout(ctx, &servicemanagementpb.CreateServiceRolloutRequest{
				ServiceName: args[0],
				Rollout: &servicemanagementpb.Rollout{
					ServiceName: args[0],
					Strategy: &servicemanagementpb.Rollout_TrafficPercentStrategy_{
						TrafficPercentStrategy: &servicemanagementpb.Rollout_TrafficPercentStrategy{
							Percentages: map[string]float64{
								args[1]: 100.0,
							},
						},
					},
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
