package servicemanagement

import (
	servicemanagement "cloud.google.com/go/servicemanagement/apiv1"
	"github.com/spf13/cobra"
)

func getServiceConfigCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "get-service-config",
		Short: "Get service config",
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
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
