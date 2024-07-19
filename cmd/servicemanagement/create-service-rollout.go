package servicemanagement

import (
	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"github.com/spf13/cobra"
)

func createServiceRolloutCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "create-service-rollout",
		Short: "Create service rollout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := servicemanagement.NewServiceManagerClient(ctx)
			if err != nil {
				return nil
			}
			defer c.Close()
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "json", "output format")
	return cmd
}
